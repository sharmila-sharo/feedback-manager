import { useEffect, useState } from 'react';

type Feedback = {
    id: number;
    message: string;
};

function App() {
    const [feedback, setFeedback] = useState<Feedback[]>([]);
    const [newMessage, setNewMessage] = useState('');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    useEffect(() => {
        setLoading(true);
        fetch(`${import.meta.env.VITE_API_URL}/feedback`)
            .then(res => {
                if (!res.ok) throw new Error("Failed to fetch feedback");
                return res.json();
            })
            .then(data => {
                setFeedback(Array.isArray(data) ? data : []);
                setError(null);
            })
            .catch(err => setError(err.message))
            .finally(() => setLoading(false));
    }, []);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setSuccess(null);

        fetch(`${import.meta.env.VITE_API_URL}/feedback`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ message: newMessage }),
        })
            .then(res => {
                if (!res.ok) throw new Error("Failed to submit feedback");
                return res.json();
            })
            .then(added => {
                setFeedback(prev => [...prev, added]);
                setNewMessage('');
                setSuccess("Feedback submitted!");
            })
            .catch(err => setError(err.message));
    };

    return (
        <div style={{ padding: '2rem' }}>
            <h1>Feedback List</h1>

            {loading && <p>Loading feedback...</p>}
            {error && <p style={{ color: 'red' }}>Error: {error}</p>}
            {success && <p style={{ color: 'green' }}>{success}</p>}

            {Array.isArray(feedback) && feedback.length > 0 ? (
                <ul>
                    {feedback.map(item => (
                        <li key={item.id}>{item.message}</li>
                    ))}
                </ul>
            ) : (
                !loading && !error && <p>No feedback yet.</p>
            )}

            <form onSubmit={handleSubmit} style={{ marginTop: '2rem' }}>
                <input
                    type="text"
                    value={newMessage}
                    onChange={e => setNewMessage(e.target.value)}
                    placeholder="Enter feedback"
                    required
                />
                <button type="submit">Submit</button>
            </form>
        </div>
    );
}

export default App;
