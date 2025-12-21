'use client';

import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import styles from './ProblemDescription.module.css';

export default function ProblemDescription({ problem }) {
    if (!problem) return null;

    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <h1 className={styles.title}>{problem.title}</h1>
                <div className={styles.limits}>
                    <span className={styles.limit}>
                        ⏱ {problem.time_limit_ms}ms
                    </span>
                    <span className={styles.limit}>
                        💾 {problem.memory_limit_mb}MB
                    </span>
                </div>
            </div>

            <div className={styles.description}>
                <ReactMarkdown remarkPlugins={[remarkGfm]}>
                    {problem.description}
                </ReactMarkdown>
            </div>

            {problem.sample_cases && problem.sample_cases.length > 0 && (
                <div className={styles.samples}>
                    <h3 className={styles.samplesTitle}>Sample Test Cases</h3>
                    {problem.sample_cases.map((tc, idx) => (
                        <div key={tc.id} className={styles.sampleCase}>
                            <div className={styles.sampleHeader}>Example {idx + 1}</div>
                            <div className={styles.sampleContent}>
                                <div className={styles.sampleBlock}>
                                    <div className={styles.sampleLabel}>Input</div>
                                    <pre className={styles.sampleCode}>{tc.input}</pre>
                                </div>
                                <div className={styles.sampleBlock}>
                                    <div className={styles.sampleLabel}>Output</div>
                                    <pre className={styles.sampleCode}>{tc.expected_output}</pre>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
