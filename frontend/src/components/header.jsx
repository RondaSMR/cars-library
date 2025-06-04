import React from "react";
import { Link, useNavigate } from "react-router-dom";

export default function Header() {
    const navigate = useNavigate();

    const handleFilter = (filter) => {
        navigate(`/?filter=${filter}`);
    };

    return (
        <header className="flex justify-between items-center px-6 py-4 bg-gray-100 shadow">
            <h1 className="text-xl font-bold">CarLibrary</h1>
            <div className="space-x-4">
                <button onClick={() => handleFilter("all")} className="text-blue-600 hover:underline">
                    Все
                </button>
                <button onClick={() => handleFilter("books")} className="text-blue-600 hover:underline">
                    Книги
                </button>
                <button onClick={() => handleFilter("drawings")} className="text-blue-600 hover:underline">
                    Схемы
                </button>
            </div>
            <Link to="/login" className="text-sm text-gray-700 hover:underline">
                Войти
            </Link>
        </header>
    );
}
