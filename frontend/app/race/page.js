'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslation } from 'react-i18next';
import { useUser } from '../hooks/useUser';
import AuthScreen from '../components/AuthScreen';
import DashboardLayout from '../components/DashboardLayout';
import { createRace, joinRace, watchRace, fetchProblems } from '../lib/api';
import styles from './page.module.css';
import { useEffect } from 'react';

export default function RaceLobby() {
    const router = useRouter();
    const { t } = useTranslation();
    const { user, loading: userLoading, login, register, logout } = useUser();
    const [mode, setMode] = useState(null); // 'create', 'join', or 'watch'
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
            setError(t('race.errors.createFailed'));
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
            setError(t('race.errors.joinFailed'));
        }
    };

    const handleWatchRoom = async () => {
        if (!roomCode.trim()) return;
        try {
            setError(null);
            await watchRace(roomCode.toUpperCase(), user.id);
            router.push(`/race/${roomCode.toUpperCase()}`);
        } catch (err) {
            setError(t('race.errors.watchFailed'));
        }
    };

    if (userLoading) {
        return <div className={styles.loading}><div className={styles.spinner}></div></div>;
    }

    if (!user) {
        return <AuthScreen onLogin={login} onRegister={register} />;
    }

    return (
        <DashboardLayout user={user} onLogout={logout}>
            <div className={styles.lobbyContainer}>
                <h1 className={styles.title}>🏁 {t('race.title')}</h1>

                {!mode && (
                    <div className={styles.options}>
                        <button className={styles.optionBtn} onClick={() => setMode('create')}>
                            <span className={styles.optionIcon}>🎮</span>
                            <span className={styles.optionText}>{t('race.createRoom')}</span>
                        </button>
                        <button className={styles.optionBtn} onClick={() => setMode('join')}>
                            <span className={styles.optionIcon}>🚀</span>
                            <span className={styles.optionText}>{t('race.joinAsPlayer')}</span>
                        </button>
                        <button className={`${styles.optionBtn} ${styles.spectatorBtn}`} onClick={() => setMode('watch')}>
                            <span className={styles.optionIcon}>👁️</span>
                            <span className={styles.optionText}>{t('race.watchAsSpectator')}</span>
                        </button>
                    </div>
                )}

                {mode === 'create' && (
                    <div className={styles.form}>
                        <h2>{t('race.createRoom')}</h2>
                        <label className={styles.label}>{t('race.selectProblem')}</label>
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
                            {creating ? t('race.creating') : t('race.createRoom')}
                        </button>
                        <button className={styles.backLink} onClick={() => setMode(null)}>
                            ← {t('common.back')}
                        </button>
                    </div>
                )}

                {mode === 'join' && (
                    <div className={styles.form}>
                        <h2>{t('race.joinAsPlayer')}</h2>
                        <p className={styles.hint}>{t('race.playerHint')}</p>
                        <label className={styles.label}>{t('race.roomCode')}</label>
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
                            {t('race.join')}
                        </button>
                        <button className={styles.backLink} onClick={() => setMode(null)}>
                            ← {t('common.back')}
                        </button>
                    </div>
                )}

                {mode === 'watch' && (
                    <div className={styles.form}>
                        <h2>{t('race.watchAsSpectator')}</h2>
                        <p className={styles.hint}>{t('race.spectatorHint')}</p>
                        <label className={styles.label}>{t('race.roomCode')}</label>
                        <input
                            type="text"
                            className={styles.input}
                            placeholder="ABCD12"
                            value={roomCode}
                            onChange={(e) => setRoomCode(e.target.value.toUpperCase())}
                            maxLength={6}
                        />
                        <button
                            className={`${styles.submitBtn} ${styles.spectatorSubmit}`}
                            onClick={handleWatchRoom}
                            disabled={!roomCode.trim()}
                        >
                            👁️ {t('race.watch')}
                        </button>
                        <button className={styles.backLink} onClick={() => setMode(null)}>
                            ← {t('common.back')}
                        </button>
                    </div>
                )}

                {error && <div className={styles.error}>{error}</div>}
            </div>
        </DashboardLayout>
    );
}
