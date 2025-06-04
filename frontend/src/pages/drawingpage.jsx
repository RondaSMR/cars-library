import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import axios from "axios";

export default function DrawingPage() {
    const { id } = useParams();
    const [drawing, setDrawing] = useState(null);
    const [comments, setComments] = useState([]);
    const [newComment, setNewComment] = useState("");

    useEffect(() => {
        const fetchData = async () => {
            try {
                const drawingRes = await axios.get(`http://localhost:8080/cars_library/drawing?id=${id}`, {
                    headers: { Authorization: "Basic " + btoa("admin:admin") },
                });
                setDrawing(drawingRes.data.data);
                setComments(drawingRes.data.data.comments || []);
            } catch (err) {
                console.error(err);
            }
        };
        fetchData();
    }, [id]);

    const handleCommentSubmit = async (e) => {
        e.preventDefault();
        // Здесь должен быть POST-запрос на создание комментария
        setComments((prev) => [
            ...prev,
            { username: "User", content: newComment, created_at: new Date().toISOString() },
        ]);
        setNewComment("");
    };

    if (!drawing) return <p className="p-6">Загрузка...</p>;

    return (
        <div className="p-6 max-w-2xl mx-auto">
            <h2 className="text-2xl font-bold mb-2">{drawing.title}</h2>
            <p className="text-sm text-gray-600 mb-4">Модель: {drawing.car_model}</p>
            <a
                href={drawing.file_url}
                className="inline-block mb-6 text-blue-600 underline"
                target="_blank"
                rel="noreferrer"
            >
                Скачать схему
            </a>

            <h3 className="font-semibold text-lg mb-2">Комментарии</h3>
            <ul className="mb-4 space-y-2">
                {comments.map((c, idx) => (
                    <li key={idx} className="bg-gray-100 p-2 rounded">
                        <p className="text-sm font-medium">{c.username}</p>
                        <p className="text-sm text-gray-700">{c.content}</p>
                        <p className="text-xs text-gray-400">{new Date(c.created_at).toLocaleString()}</p>
                    </li>
                ))}
            </ul>

            <form onSubmit={handleCommentSubmit} className="space-y-2">
        <textarea
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
            className="w-full border p-2 rounded"
            rows={3}
            placeholder="Оставьте комментарий..."
            required
        ></textarea>
                <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
                    Отправить
                </button>
            </form>
        </div>
    );
}
