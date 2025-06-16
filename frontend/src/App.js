import "./App.css";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import HomePage from "./Pages/HomePage";
import GamePage from "./Pages/GamePage";
import ModelsPage from "./Pages/ModelsPage";

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
