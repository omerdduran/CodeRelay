'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useUser } from '../hooks/useUser';
import { createRace, joinRace, fetchProblems } from '../lib/api';
import NicknameScreen from '../components/NicknameScreen';
import styles from './page.module.css';
import Link from 'next/link';
import { useEffect } from 'react';

export default function RaceLobby() {
    const router = useRouter();
    const { user, loading: userLoading, login } = useUser();
    const [mode, setMode] = useState(null); // 'create' or 'join'
    const [roomCode, setRoomCode] = useState('');
    const [problems, setProblems] = useState([]);
    const [selectedProblem, setSelectedProblem] = useState(1);
    const [error, setError] = useState(null);
    const [creating, setCreating] = useState(false);

    useEffect(() => {
        fetchProblems().then(setProblems).catch(() => { });
    }, []);

    const handleCreateRoom = async () => {
        try {
            setCreating(true);
            setError(null);
            const race = await createRace(user.id, selectedProblem);
            router.push(`/race/${race.room_code}`);
        } catch (err) {
            setError('Failed to create room');
        } finally {
            setCreating(false);
        }
    };

    const handleJoinRoom = async () => {
        if (!roomCode.trim()) return;
        try {
            setError(null);
            await joinRace(roomCode.toUpperCase(), user.id);
            router.push(`/race/${roomCode.toUpperCase()}`);
        } catch (err) {
            setError('Room not found or already started');
        }
    };

    if (userLoading) {
        return <div className={styles.loading}><div className={styles.spinner}></div></div>;
    }

    if (!user) {
        return <NicknameScreen onNicknameSet={login} />;
    }

    return (
        <div className={styles.container}>
            <header className={styles.header}>
                <Link href="/" className={styles.backBtn}>← Home</Link>
                <span className={styles.nickname}>{user.nickname}</span>
            </header>

            <main className={styles.main}>
                <h1 className={styles.title}>⚡ Race Mode</h1>

                {!mode && (
                    <div className={styles.options}>
                        <button className={styles.optionBtn} onClick={() => setMode('create')}>
                            🏁 Create Room
                        </button>
                        <button className={styles.optionBtn} onClick={() => setMode('join')}>
                            🚀 Join Room
                        </button>
                    </div>
                )}

                {mode === 'create' && (
                    <div className={styles.form}>
                        <h2>Create a Room</h2>
                        <label className={styles.label}>Select Problem</label>
                        <select
                            className={styles.select}
                            value={selectedProblem}
                            onChange={(e) => setSelectedProblem(Number(e.target.value))}
                        >
                            {problems.map(p => (
                                <option key={p.id} value={p.id}>{p.title}</option>
                            ))}
                        </select>
                        <button
                            className={styles.submitBtn}
                            onClick={handleCreateRoom}
                            disabled={creating}
                        >
                            {creating ? 'Creating...' : 'Create Room'}
                        </button>
                        <button className={styles.backLink} onClick={() => setMode(null)}>
                            ← Back
                        </button>
                    </div>
                )}

                {mode === 'join' && (
                    <div className={styles.form}>
                        <h2>Join a Room</h2>
                        <label className={styles.label}>Room Code</label>
                        <input
                            type="text"
                            className={styles.input}
                            placeholder="ABCD12"
                            value={roomCode}
                            onChange={(e) => setRoomCode(e.target.value.toUpperCase())}
                            maxLength={6}
                        />
                        <button
                            className={styles.submitBtn}
                            onClick={handleJoinRoom}
                            disabled={!roomCode.trim()}
                        >
                            Join Room
                        </button>
                        <button className={styles.backLink} onClick={() => setMode(null)}>
                            ← Back
                        </button>
                    </div>
                )}

                {error && <div className={styles.error}>{error}</div>}
            </main>
        </div>
    );
}
