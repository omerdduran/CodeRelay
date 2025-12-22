'use client';

import { useEffect, useState, useCallback, use } from 'react';
import { useRouter } from 'next/navigation';
import { useUser } from '../../hooks/useUser';
import { useWebSocket } from '../../hooks/useWebSocket';
import { fetchRace, startRace, createSubmission, fetchSubmission, fetchProblem } from '../../lib/api';
import NicknameScreen from '../../components/NicknameScreen';
import CodeEditor from '../../components/CodeEditor';
import ProblemDescription from '../../components/ProblemDescription';
import styles from './page.module.css';
import Link from 'next/link';

export default function RaceRoom({ params }) {
    const { code } = use(params);
    const router = useRouter();
    const { user, loading: userLoading, login } = useUser();

    const [race, setRace] = useState(null);
    const [problem, setProblem] = useState(null);
    const [loading, setLoading] = useState(true);
    const [countdown, setCountdown] = useState(null);
    const [raceTime, setRaceTime] = useState(0);
    const [codeText, setCodeText] = useState('');
    const [submitting, setSubmitting] = useState(false);
    const [myVerdict, setMyVerdict] = useState(null);

    // Handle WebSocket messages
    const handleMessage = useCallback((message) => {
        if (message.type === 'race_event') {
            // payload might be string or already parsed object
            const payload = typeof message.payload === 'string'
                ? JSON.parse(message.payload)
                : message.payload;

            if (payload.room_code !== code) return;

            if (payload.event === 'player_joined') {
                loadRace();
            } else if (payload.event === 'countdown') {
                setCountdown(payload.seconds);
                const interval = setInterval(() => {
                    setCountdown(prev => {
                        if (prev <= 1) {
                            clearInterval(interval);
                            return null;
                        }
                        return prev - 1;
                    });
                }, 1000);
            } else if (payload.event === 'race_started') {
                setRace(prev => prev ? { ...prev, status: 'racing' } : prev);
                // Start timer
                const startTime = new Date(payload.start_time).getTime();
                const timerInterval = setInterval(() => {
                    setRaceTime(Math.floor((Date.now() - startTime) / 1000));
                }, 100);
                // Clear on unmount
                return () => clearInterval(timerInterval);
            } else if (payload.event === 'race_progress') {
                loadRace();
            }
        }
    }, [code]);

    const { connected } = useWebSocket(handleMessage);

    const loadRace = async () => {
        try {
            const data = await fetchRace(code);
            setRace(data);
            if (data.problem_id) {
                const prob = await fetchProblem(data.problem_id);
                setProblem(prob);
            }
        } catch (err) {
            router.push('/race');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (user && code) {
            loadRace();
        }
    }, [user, code]);

    const handleStart = async () => {
        try {
            await startRace(code, user.id);
        } catch (err) {
            console.error('Failed to start race:', err);
        }
    };

    const handleSubmit = async () => {
        if (!codeText.trim() || submitting) return;

        try {
            setSubmitting(true);
            const result = await createSubmission(user.id, race.problem_id, codeText);

            // Poll for result
            const poll = setInterval(async () => {
                const updated = await fetchSubmission(result.id);
                if (updated.status !== 'queued' && updated.status !== 'running') {
                    clearInterval(poll);
                    setMyVerdict(updated.status);
                    setSubmitting(false);
                }
            }, 500);
        } catch (err) {
            setSubmitting(false);
        }
    };

    if (userLoading || loading) {
        return <div className={styles.loading}><div className={styles.spinner}></div></div>;
    }

    if (!user) {
        return <NicknameScreen onNicknameSet={login} />;
    }

    if (!race) {
        return <div className={styles.error}>Race not found</div>;
    }

    // Debug: log IDs for troubleshooting
    console.log('User ID:', user.id, 'Host ID:', race.host_user_id);

    const isHost = Number(race.host_user_id) === Number(user.id);
    const isWaiting = race.status === 'waiting';
    const isRacing = race.status === 'racing' || countdown !== null;

    // Waiting Room
    if (isWaiting && countdown === null) {
        return (
            <div className={styles.container}>
                <header className={styles.header}>
                    <Link href="/race" className={styles.backBtn}>← Leave</Link>
                    <span className={styles.roomCode}>Room: {code}</span>
                </header>

                <main className={styles.waitingRoom}>
                    <h1 className={styles.title}>Waiting for players...</h1>

                    <div className={styles.code}>
                        <span className={styles.codeLabel}>Share this code</span>
                        <span className={styles.codeValue}>{code}</span>
                    </div>

                    <div className={styles.players}>
                        <h3>Players ({race.players?.length || 0})</h3>
                        <ul>
                            {race.players?.map(p => (
                                <li key={p.user_id} className={styles.player}>
                                    {p.nickname}
                                    {p.user_id === race.host_user_id && <span className={styles.hostBadge}>Host</span>}
                                </li>
                            ))}
                        </ul>
                    </div>

                    {isHost && (
                        <button
                            className={styles.startBtn}
                            onClick={handleStart}
                            disabled={race.players?.length < 1}
                        >
                            🏁 Start Race
                        </button>
                    )}

                    {!isHost && (
                        <div className={styles.waitingMsg}>Waiting for host to start...</div>
                    )}
                </main>
            </div>
        );
    }

    // Countdown
    if (countdown !== null) {
        return (
            <div className={styles.countdownScreen}>
                <div className={styles.countdownNumber}>{countdown}</div>
                <div className={styles.countdownText}>Get Ready!</div>
            </div>
        );
    }

    // Check if race is finished (all players have verdict or user got AC)
    const allFinished = race.players?.every(p => p.verdict);
    const showResults = allFinished || myVerdict === 'AC';

    // Results View
    if (showResults && myVerdict) {
        const sortedPlayers = [...(race.players || [])].sort((a, b) => {
            if (a.verdict === 'AC' && b.verdict !== 'AC') return -1;
            if (a.verdict !== 'AC' && b.verdict === 'AC') return 1;
            return (a.finish_time || 999999) - (b.finish_time || 999999);
        });

        return (
            <div className={styles.resultsScreen}>
                <div className={styles.resultsCard}>
                    <h1 className={styles.resultsTitle}>🏁 Race Results</h1>

                    <div className={styles.resultsList}>
                        {sortedPlayers.map((p, idx) => (
                            <div
                                key={p.user_id}
                                className={`${styles.resultItem} ${p.user_id === user.id ? styles.me : ''}`}
                            >
                                <span className={styles.resultRank}>
                                    {idx === 0 ? '🥇' : idx === 1 ? '🥈' : idx === 2 ? '🥉' : `#${idx + 1}`}
                                </span>
                                <span className={styles.resultName}>{p.nickname}</span>
                                <span className={`${styles.resultVerdict} ${p.verdict === 'AC' ? styles.ac : styles.wa}`}>
                                    {p.verdict || '-'}
                                </span>
                                <span className={styles.resultTime}>
                                    {p.finish_time ? `${(p.finish_time / 1000).toFixed(1)}s` : '-'}
                                </span>
                            </div>
                        ))}
                    </div>

                    <Link href="/race" className={styles.newRaceBtn}>
                        🔄 New Race
                    </Link>
                </div>
            </div>
        );
    }

    // Racing View
    return (
        <div className={styles.raceContainer}>
            <header className={styles.raceHeader}>
                <span className={styles.roomCode}>Race: {code}</span>
                <span className={styles.timer}>⏱ {Math.floor(raceTime / 60)}:{String(raceTime % 60).padStart(2, '0')}</span>
                <span className={styles.wsStatus}>{connected ? '🟢' : '🔴'}</span>
            </header>

            <div className={styles.raceContent}>
                <div className={styles.leftPanel}>
                    <ProblemDescription problem={problem} />
                </div>

                <div className={styles.rightPanel}>
                    <div className={styles.editorSection}>
                        <CodeEditor onCodeChange={setCodeText} />
                    </div>

                    <div className={styles.actionBar}>
                        {myVerdict ? (
                            <div className={`${styles.verdict} ${myVerdict === 'AC' ? styles.ac : styles.wa}`}>
                                {myVerdict === 'AC' ? '✓ Accepted!' : `✗ ${myVerdict}`}
                            </div>
                        ) : (
                            <button
                                className={styles.submitBtn}
                                onClick={handleSubmit}
                                disabled={submitting || !codeText.trim()}
                            >
                                {submitting ? 'Submitting...' : '▶ Submit'}
                            </button>
                        )}
                    </div>

                    <div className={styles.playersBar}>
                        {race.players?.map(p => (
                            <div key={p.user_id} className={styles.playerStatus}>
                                <span>{p.nickname}</span>
                                <span className={styles.playerVerdict}>
                                    {p.verdict === 'AC' ? '✓' : p.status === 'racing' ? '⌛' : '-'}
                                </span>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
}

