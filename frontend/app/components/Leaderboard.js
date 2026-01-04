'use client';

import { useEffect, useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { useWebSocket } from '../hooks/useWebSocket';
import styles from './Leaderboard.module.css';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export default function Leaderboard({ problemId = 1 }) {
    const { t } = useTranslation();
    const [entries, setEntries] = useState([]);
    const [loading, setLoading] = useState(true);

    const fetchLeaderboard = useCallback(async () => {
        try {
            const res = await fetch(`${API_URL}/api/leaderboard?problem_id=${problemId}`);
            if (res.ok) {
                const data = await res.json();
                setEntries(data || []);
            }
        } catch (err) {
            console.error('Failed to fetch leaderboard:', err);
        } finally {
            setLoading(false);
        }
    }, [problemId]);

    // Handle WebSocket messages
    const handleMessage = useCallback((message) => {
        if (message.type === 'submission_update') {
            const payload = JSON.parse(message.payload);
            if (payload.status === 'AC' && payload.problem_id === problemId) {
                // Refresh leaderboard when someone solves
                fetchLeaderboard();
            }
        }
    }, [problemId, fetchLeaderboard]);

    const { connected } = useWebSocket(handleMessage);

    useEffect(() => {
        fetchLeaderboard();
    }, [fetchLeaderboard]);

    if (loading) {
        return (
            <div className={styles.container}>
                <div className={styles.loading}>{t('leaderboard.loading')}</div>
            </div>
        );
    }

    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <h3 className={styles.title}>{t('leaderboard.title')}</h3>
                <span className={styles.status}>
                    {connected ? t('leaderboard.live') : t('leaderboard.offline')}
                </span>
            </div>

            {entries.length === 0 ? (
                <div className={styles.empty}>{t('leaderboard.empty')}</div>
            ) : (
                <div className={styles.list}>
                    {entries.map((entry) => (
                        <div key={entry.user_id} className={styles.entry}>
                            <span className={styles.rank}>#{entry.rank}</span>
                            <span className={styles.nickname}>{entry.nickname}</span>
                            <span className={styles.time}>
                                {entry.solve_time_ms ? `${entry.solve_time_ms}ms` : '-'}
                            </span>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
