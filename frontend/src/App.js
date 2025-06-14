import "./App.css";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import HomePage from "./pages/HomePage";
import GamePage from "./pages/GamePage";
import ModelsPage from "./pages/ModelsPage";

import Modal from "react-modal";

Modal.setAppElement("#root");

const App = () => {
    return (
        <Router>
            <Routes>
                <Route path="/" element={<HomePage />} />
                <Route path="/game" element={<GamePage />} />
                <Route path="/models" element={<ModelsPage />} />
            </Routes>
        </Router>
    );
};

export default App;
