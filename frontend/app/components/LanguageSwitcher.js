'use client';

import { useTranslation } from 'react-i18next';
import styles from './LanguageSwitcher.module.css';

const LANGUAGES = [
    { code: 'en', name: 'English', flag: '🇬🇧' },
    { code: 'tr', name: 'Türkçe', flag: '🇹🇷' },
    { code: 'pt', name: 'Português', flag: '🇧🇷' }
];

export default function LanguageSwitcher() {
    const { i18n } = useTranslation();

    const handleChange = (e) => {
        i18n.changeLanguage(e.target.value);
    };

    return (
        <div className={styles.container}>
            <select
                className={styles.select}
                value={i18n.language}
                onChange={handleChange}
            >
                {LANGUAGES.map((lang) => (
                    <option key={lang.code} value={lang.code}>
                        {lang.flag} {lang.name}
                    </option>
                ))}
            </select>
        </div>
    );
}
