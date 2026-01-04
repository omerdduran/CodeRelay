'use client';

import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import LanguageSwitcher from './LanguageSwitcher';
import styles from './NicknameScreen.module.css';

export default function NicknameScreen({ onNicknameSet }) {
    const { t } = useTranslation();
    const [nickname, setNickname] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();

        const trimmed = nickname.trim();
        if (!trimmed) {
            setError(t('nickname.errors.required'));
            return;
        }
        if (trimmed.length < 2) {
            setError(t('nickname.errors.tooShort'));
            return;
        }
        if (trimmed.length > 20) {
            setError(t('nickname.errors.tooLong'));
            return;
        }

        setLoading(true);
        setError('');

        try {
            await onNicknameSet(trimmed);
        } catch (err) {
            setError(t('nickname.errors.failed'));
            setLoading(false);
        }
    };

    return (
        <div className={styles.container}>
            <div className={styles.languageSwitcherWrapper}>
                <LanguageSwitcher />
            </div>
            <div className={styles.card}>
                <div className={styles.logo}>
                    <span className={styles.logoIcon}>⚡</span>
                    <h1 className={styles.title}>CodeRelay</h1>
                </div>
                <p className={styles.subtitle}>{t('nickname.subtitle')}</p>

                <form onSubmit={handleSubmit} className={styles.form}>
                    <input
                        type="text"
                        value={nickname}
                        onChange={(e) => setNickname(e.target.value)}
                        placeholder={t('nickname.placeholder')}
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
                        {loading ? t('nickname.joining') : t('nickname.joinButton')}
                    </button>
                </form>
            </div>
        </div>
    );
}
