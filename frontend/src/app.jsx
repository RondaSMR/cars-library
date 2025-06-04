import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import HomePage from "./pages/homepage";
import BookPage from "./pages/bookpage";
import DrawingPage from "./pages/drawingpage";
import Header from "./components/header";

export default function App() {
    return (
        <Router>
            <Header />
            <Routes>
                <Route path="/" element={<HomePage />} />
                <Route path="/book/:id" element={<BookPage />} />
                <Route path="/drawing/:id" element={<DrawingPage />} />
            </Routes>
        </Router>
    );
}
