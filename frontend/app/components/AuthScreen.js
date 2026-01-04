'use client';

import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import LanguageSwitcher from './LanguageSwitcher';
import styles from './AuthScreen.module.css';

export default function AuthScreen({ onLogin, onRegister }) {
    const { t } = useTranslation();
    const [mode, setMode] = useState('login'); // 'login' or 'register'
    const [nickname, setNickname] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const validateForm = () => {
        const trimmedNickname = nickname.trim();

        if (!trimmedNickname) {
            setError(t('auth.errors.nicknameRequired'));
            return false;
        }
        if (trimmedNickname.length < 2) {
            setError(t('auth.errors.nicknameTooShort'));
            return false;
        }
        if (trimmedNickname.length > 20) {
            setError(t('auth.errors.nicknameTooLong'));
            return false;
        }
        if (!password) {
            setError(t('auth.errors.passwordRequired'));
            return false;
        }
        if (password.length < 6) {
            setError(t('auth.errors.passwordTooShort'));
            return false;
        }
        if (mode === 'register' && password !== confirmPassword) {
            setError(t('auth.errors.passwordMismatch'));
            return false;
        }

        return true;
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        if (!validateForm()) return;

        setLoading(true);

        try {
            if (mode === 'login') {
                await onLogin(nickname.trim(), password);
            } else {
                await onRegister(nickname.trim(), password);
            }
        } catch (err) {
            if (err.message === 'nickname_taken') {
                setError(t('auth.errors.nicknameTaken'));
            } else if (err.message === 'invalid_credentials') {
                setError(t('auth.errors.invalidCredentials'));
            } else if (mode === 'login') {
                setError(t('auth.errors.loginFailed'));
            } else {
                setError(t('auth.errors.registrationFailed'));
            }
        } finally {
            setLoading(false);
        }
    };

    const toggleMode = () => {
        setMode(mode === 'login' ? 'register' : 'login');
        setError('');
        setConfirmPassword('');
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

                <div className={styles.tabs}>
                    <button
                        type="button"
                        className={`${styles.tab} ${mode === 'login' ? styles.activeTab : ''}`}
                        onClick={() => { setMode('login'); setError(''); }}
                    >
                        {t('auth.login')}
                    </button>
                    <button
                        type="button"
                        className={`${styles.tab} ${mode === 'register' ? styles.activeTab : ''}`}
                        onClick={() => { setMode('register'); setError(''); }}
                    >
                        {t('auth.register')}
                    </button>
                </div>

                <form onSubmit={handleSubmit} className={styles.form}>
                    <input
                        type="text"
                        value={nickname}
                        onChange={(e) => setNickname(e.target.value)}
                        placeholder={t('auth.nickname')}
                        className={styles.input}
                        maxLength={20}
                        autoFocus
                        disabled={loading}
                    />
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        placeholder={t('auth.password')}
                        className={styles.input}
                        disabled={loading}
                    />
                    {mode === 'register' && (
                        <input
                            type="password"
                            value={confirmPassword}
                            onChange={(e) => setConfirmPassword(e.target.value)}
                            placeholder={t('auth.confirmPassword')}
                            className={styles.input}
                            disabled={loading}
                        />
                    )}
                    {error && <p className={styles.error}>{error}</p>}
                    <button
                        type="submit"
                        className={styles.button}
                        disabled={loading}
                    >
                        {loading
                            ? (mode === 'login' ? t('auth.loggingIn') : t('auth.registering'))
                            : (mode === 'login' ? t('auth.login') : t('auth.register'))
                        }
                    </button>
                </form>

                <p className={styles.switchText}>
                    {mode === 'login' ? t('auth.noAccount') : t('auth.hasAccount')}{' '}
                    <button type="button" className={styles.switchLink} onClick={toggleMode}>
                        {mode === 'login' ? t('auth.register') : t('auth.login')}
                    </button>
                </p>
            </div>
        </div>
    );
}
