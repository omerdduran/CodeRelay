'use client';

import { useEffect, useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { useWebSocket } from '../hooks/useWebSocket';
import styles from './Leaderboard.module.css';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export default function Leaderboard({ problemId = null, showElo = true }) {
    const { t } = useTranslation();
    const [entries, setEntries] = useState([]);
    const [loading, setLoading] = useState(true);
    const [viewMode, setViewMode] = useState(showElo && !problemId ? 'elo' : 'problem');

    const fetchLeaderboard = useCallback(async () => {
        try {
            let url = `${API_URL}/api/leaderboard`;
            
            if (viewMode === 'elo') {
                url += '?type=elo&limit=100';
            } else if (problemId) {
                url += `?problem_id=${problemId}`;
            } else {
                url += '?type=elo&limit=100';
            }
            
            const res = await fetch(url);
            if (res.ok) {
                const data = await res.json();
                setEntries(data || []);
            }
        } catch (err) {
            console.error('Failed to fetch leaderboard:', err);
        } finally {
            setLoading(false);
        }
    }, [problemId, viewMode]);

    // Handle WebSocket messages
    const handleMessage = useCallback((message) => {
        if (message.type === 'submission_update') {
            const payload = JSON.parse(message.payload);
            if (payload.status === 'AC') {
                // Refresh leaderboard when someone solves
                fetchLeaderboard();
            }
        } else if (message.type === 'race_event') {
            const payload = JSON.parse(message.payload);
            if (payload.event === 'elo_updated') {
                // Refresh ELO leaderboard when ratings are updated
                if (viewMode === 'elo') {
                    fetchLeaderboard();
                }
            }
        }
    }, [viewMode, fetchLeaderboard]);

    const { connected } = useWebSocket(handleMessage);

    useEffect(() => {
        fetchLeaderboard();
    }, [fetchLeaderboard]);

    const toggleViewMode = () => {
        setViewMode(viewMode === 'elo' ? 'problem' : 'elo');
    };

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
                <h3 className={styles.title}>
                    {viewMode === 'elo' ? t('leaderboard.globalRanking') : t('leaderboard.problemRanking')}
                </h3>
                <span className={styles.status}>
                    {connected ? t('leaderboard.live') : t('leaderboard.offline')}
                </span>
            </div>

            {showElo && !problemId && (
                <div className={styles.toggleContainer}>
                    <button 
                        className={`${styles.toggleButton} ${viewMode === 'elo' ? styles.active : ''}`}
                        onClick={() => setViewMode('elo')}
                    >
                        {t('leaderboard.viewElo')}
                    </button>
                    <button 
                        className={`${styles.toggleButton} ${viewMode === 'problem' ? styles.active : ''}`}
                        onClick={() => setViewMode('problem')}
                    >
                        {t('leaderboard.viewProblem')}
                    </button>
                </div>
            )}

            {entries.length === 0 ? (
                <div className={styles.empty}>{t('leaderboard.empty')}</div>
            ) : (
                <div className={styles.list}>
                    {entries.map((entry) => (
                        <div key={entry.user_id} className={styles.entry}>
                            <span className={styles.rank}>#{entry.rank}</span>
                            <span className={styles.nickname}>{entry.nickname}</span>
                            {viewMode === 'elo' ? (
                                <span className={styles.elo}>
                                    {entry.elo_rating || 1200}
                                </span>
                            ) : (
                                <span className={styles.time}>
                                    {entry.solve_time_ms ? `${entry.solve_time_ms}ms` : '-'}
                                </span>
                            )}
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
