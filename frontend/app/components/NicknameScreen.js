'use client';

import { useState } from 'react';
import styles from './NicknameScreen.module.css';

export default function NicknameScreen({ onNicknameSet }) {
    const [nickname, setNickname] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();

        const trimmed = nickname.trim();
        if (!trimmed) {
            setError('Please enter a nickname');
            return;
        }
        if (trimmed.length < 2) {
            setError('Nickname must be at least 2 characters');
            return;
        }
        if (trimmed.length > 20) {
            setError('Nickname must be 20 characters or less');
            return;
        }

        setLoading(true);
        setError('');

        try {
            await onNicknameSet(trimmed);
        } catch (err) {
            setError('Failed to set nickname. Please try again.');
            setLoading(false);
        }
    };

    return (
        <div className={styles.container}>
            <div className={styles.card}>
                <div className={styles.logo}>
                    <span className={styles.logoIcon}>⚡</span>
                    <h1 className={styles.title}>CodeRelay</h1>
                </div>
                <p className={styles.subtitle}>Real-time competitive coding</p>

                <form onSubmit={handleSubmit} className={styles.form}>
                    <input
                        type="text"
                        value={nickname}
                        onChange={(e) => setNickname(e.target.value)}
                        placeholder="Enter your nickname"
                        className={styles.input}
                        maxLength={20}
                        autoFocus
                        disabled={loading}
                    />
                    {error && <p className={styles.error}>{error}</p>}
                    <button
                        type="submit"
                        className={styles.button}
                        disabled={loading}
                    >
                        {loading ? 'Joining...' : 'Join Arena'}
                    </button>
                </form>
            </div>
        </div>
    );
}
