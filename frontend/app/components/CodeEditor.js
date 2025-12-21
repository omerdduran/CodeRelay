'use client';

import { useState } from 'react';
import Editor from '@monaco-editor/react';
import styles from './CodeEditor.module.css';

const DEFAULT_CODE = `# Read input
nums = list(map(int, input().split()))
target = int(input())

# Your solution here
def two_sum(nums, target):
    seen = {}
    for i, num in enumerate(nums):
        complement = target - num
        if complement in seen:
            return [seen[complement], i]
        seen[num] = i
    return []

# Output result
result = two_sum(nums, target)
print(result[0], result[1])
`;

export default function CodeEditor({
    initialCode = DEFAULT_CODE,
    language = 'python',
    onCodeChange,
    readOnly = false
}) {
    const [code, setCode] = useState(initialCode);

    const handleEditorChange = (value) => {
        setCode(value || '');
        if (onCodeChange) {
            onCodeChange(value || '');
        }
    };

    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <span className={styles.language}>{language}</span>
                <div className={styles.actions}>
                    <button
                        className={styles.resetBtn}
                        onClick={() => handleEditorChange(DEFAULT_CODE)}
                        title="Reset to template"
                    >
                        ↺ Reset
                    </button>
                </div>
            </div>
            <div className={styles.editorWrapper}>
                <Editor
                    height="100%"
                    language={language}
                    value={code}
                    onChange={handleEditorChange}
                    theme="vs-dark"
                    options={{
                        minimap: { enabled: false },
                        fontSize: 14,
                        fontFamily: "'Fira Code', 'Cascadia Code', Consolas, monospace",
                        lineNumbers: 'on',
                        scrollBeyondLastLine: false,
                        automaticLayout: true,
                        tabSize: 4,
                        readOnly,
                        padding: { top: 16 },
                    }}
                />
            </div>
        </div>
    );
}
