'use client';

import styles from './SubmitButton.module.css';

export default function SubmitButton({ onSubmit, loading, disabled }) {
    return (
        <button
            className={styles.button}
            onClick={onSubmit}
            disabled={loading || disabled}
        >
            {loading ? (
                <>
                    <span className={styles.spinner}></span>
                    Submitting...
                </>
            ) : (
                <>
                    <span className={styles.icon}>▶</span>
                    Submit
                </>
            )}
        </button>
    );
}
