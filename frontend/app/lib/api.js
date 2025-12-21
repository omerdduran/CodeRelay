const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export async function fetchProblems() {
    const res = await fetch(`${API_URL}/api/problems`);
    if (!res.ok) throw new Error('Failed to fetch problems');
    return res.json();
}

export async function fetchProblem(id) {
    const res = await fetch(`${API_URL}/api/problems/${id}`);
    if (!res.ok) throw new Error('Failed to fetch problem');
    return res.json();
}

export async function createSubmission(userId, problemId, code, language = 'python') {
    const res = await fetch(`${API_URL}/api/submissions`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            user_id: userId,
            problem_id: problemId,
            code,
            language,
        }),
    });
    if (!res.ok) throw new Error('Failed to create submission');
    return res.json();
}

export async function fetchSubmission(id) {
    const res = await fetch(`${API_URL}/api/submissions/${id}`);
    if (!res.ok) throw new Error('Failed to fetch submission');
    return res.json();
}

export async function createUser(nickname) {
    const res = await fetch(`${API_URL}/api/users`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ nickname }),
    });
    if (!res.ok) throw new Error('Failed to create user');
    return res.json();
}
