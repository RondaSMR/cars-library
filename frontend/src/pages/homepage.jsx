import React, { useEffect, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import axios from "axios";

export default function HomePage() {
    const [items, setItems] = useState([]);
    const [filter, setFilter] = useState("all");
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();

    useEffect(() => {
        const f = searchParams.get("filter") || "all";
        setFilter(f);
    }, [searchParams]);

    useEffect(() => {
        const fetchItems = async () => {
            try {
                const [booksRes, drawingsRes] = await Promise.all([
                    axios.get("http://localhost:8080/cars_library/book", {
                        headers: { Authorization: "Basic " + btoa("admin:admin") },
                    }),
                    axios.get("http://localhost:8080/cars_library/drawing", {
                        headers: { Authorization: "Basic " + btoa("admin:admin") },
                    }),
                ]);
                setItems([...booksRes.data.data, ...drawingsRes.data.data]);
            } catch (err) {
                console.error(err);
            }
        };
        fetchItems();
    }, []);

    const filteredItems = items.filter((item) => {
        if (filter === "books") return item.type === "book";
        if (filter === "drawings") return item.type === "drawing";
        return true;
    });

    return (
        <div className="p-6">
            <h2 className="text-lg font-semibold mb-4">Каталог</h2>
            <div className="grid gap-4 sm:grid-cols-2 md:grid-cols-3">
                {filteredItems.map((item) => (
                    <div
                        key={item.id}
                        className="border p-4 rounded shadow hover:shadow-md cursor-pointer"
                        onClick={() => navigate(`/${item.type}/${item.id}`)}
                    >
                        <h3 className="font-bold text-lg mb-2">{item.title}</h3>
                        <p className="text-sm text-gray-600">{item.type === "book" ? "Книга" : "Схема"}</p>
                    </div>
                ))}
            </div>
        </div>
    );
}
