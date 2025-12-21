'use client';

import { useEffect, useState, use } from 'react';
import { useRouter } from 'next/navigation';
import { useUser } from '../../hooks/useUser';
import { fetchProblem, createSubmission, fetchSubmission } from '../../lib/api';
import ProblemDescription from '../../components/ProblemDescription';
import CodeEditor from '../../components/CodeEditor';
import SubmitButton from '../../components/SubmitButton';
import ResultsPanel from '../../components/ResultsPanel';
import NicknameScreen from '../../components/NicknameScreen';
import styles from './page.module.css';
import Link from 'next/link';

export default function ProblemPage({ params }) {
    const { id } = use(params);
    const router = useRouter();
    const { user, loading: userLoading, login } = useUser();

    const [problem, setProblem] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [code, setCode] = useState('');
    const [submitting, setSubmitting] = useState(false);
    const [submission, setSubmission] = useState(null);

    useEffect(() => {
        if (user && id) {
            loadProblem();
        }
    }, [user, id]);

    const loadProblem = async () => {
        try {
            setLoading(true);
            const data = await fetchProblem(id);
            setProblem(data);
        } catch (err) {
            setError('Failed to load problem');
        } finally {
            setLoading(false);
        }
    };

    const handleSubmit = async () => {
        if (!code.trim()) return;

        try {
            setSubmitting(true);
            setSubmission(null);

            const result = await createSubmission(user.id, parseInt(id), code);
            setSubmission(result);

            // Poll for result (since we don't have WebSocket yet)
            let attempts = 0;
            const pollInterval = setInterval(async () => {
                attempts++;
                try {
                    const updated = await fetchSubmission(result.id);
                    setSubmission(updated);
                    if (updated.status !== 'queued' && updated.status !== 'running') {
                        clearInterval(pollInterval);
                    }
                } catch {
                    // ignore poll errors
                }
                if (attempts >= 30) {
                    clearInterval(pollInterval);
                }
            }, 1000);
        } catch (err) {
            setError('Failed to submit code');
        } finally {
            setSubmitting(false);
        }
    };

    if (userLoading) {
        return (
            <div className={styles.loadingScreen}>
                <div className={styles.spinner}></div>
            </div>
        );
    }

    if (!user) {
        return <NicknameScreen onNicknameSet={login} />;
    }

    if (loading) {
        return (
            <div className={styles.loadingScreen}>
                <div className={styles.spinner}></div>
                <span>Loading problem...</span>
            </div>
        );
    }

    if (error) {
        return (
            <div className={styles.errorScreen}>
                <p>{error}</p>
                <Link href="/" className={styles.backLink}>← Back to problems</Link>
            </div>
        );
    }

    return (
        <div className={styles.container}>
            <header className={styles.header}>
                <Link href="/" className={styles.backBtn}>
                    ← Problems
                </Link>
                <div className={styles.userInfo}>
                    <span className={styles.nickname}>{user.nickname}</span>
                </div>
            </header>

            <div className={styles.content}>
                <div className={styles.leftPanel}>
                    <ProblemDescription problem={problem} />
                </div>

                <div className={styles.rightPanel}>
                    <div className={styles.editorSection}>
                        <CodeEditor
                            onCodeChange={setCode}
                        />
                    </div>

                    <div className={styles.actionBar}>
                        <SubmitButton
                            onSubmit={handleSubmit}
                            loading={submitting}
                            disabled={!code.trim()}
                        />
                    </div>

                    <div className={styles.resultsSection}>
                        <ResultsPanel
                            submission={submission}
                            loading={submitting}
                        />
                    </div>
                </div>
            </div>
        </div>
    );
}
