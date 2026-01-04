'use client';

import { useTranslation } from 'react-i18next';
import styles from './SubmitButton.module.css';

export default function SubmitButton({ onSubmit, loading, disabled }) {
    const { t } = useTranslation();

    return (
        <button
            className={styles.button}
            onClick={onSubmit}
            disabled={loading || disabled}
        >
            {loading ? (
                <>
                    <span className={styles.spinner}></span>
                    {t('common.submitting')}
                </>
            ) : (
                <>
                    <span className={styles.icon}>▶</span>
                    {t('common.submit')}
                </>
            )}
        </button>
    );
}
