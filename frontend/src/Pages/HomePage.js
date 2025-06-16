import React from "react";
import { useNavigate } from "react-router-dom";
import "./HomePage.css";
import Header from "../Components/Header";
import Footer from "../Components/Footer";

const HomePage = () => {
    const navigate = useNavigate();

    const handleTimeControlClick = (baseTime, increment) => {
        navigate(`/game?baseTime=${baseTime}&increment=${increment}`);
    };

    return (
        <div className="page-container">
            <Header />
            <div className="content-container">
                <div className="white-wrapper">
                    <h1>Кинь виклик шаховому рушію!</h1>
                    <p>
                        Киньте собі виклик та перевірте свої вміння проти
                        шахового рушія на основі дерева пошуку Монте-Карло та
                        згорткової нейромережі.
                    </p>
                </div>
                <div className="white-wrapper">
                    <h2>Виберіть контроль часу</h2>
                    <div className="time-buttons">
                        <button
                            className="time-button"
                            onClick={() => handleTimeControlClick(180, 2)}
                        >
                            3+2
                            <span className="time-type">
                                блискавичний контроль
                            </span>
                        </button>
                        <button
                            className="time-button"
                            onClick={() => handleTimeControlClick(300, 0)}
                        >
                            5+0
                            <span className="time-type">
                                блискавичний контроль
                            </span>
                        </button>
                        <button
                            className="time-button"
                            onClick={() => handleTimeControlClick(300, 3)}
                        >
                            5+3
                            <span className="time-type">
                                блискавичний контроль
                            </span>
                        </button>
                        <button
                            className="time-button"
                            onClick={() => handleTimeControlClick(600, 0)}
                        >
                            10+0
                            <span className="time-type">швидкий контроль</span>
                        </button>
                        <button
                            className="time-button"
                            onClick={() => handleTimeControlClick(600, 10)}
                        >
                            10+5
                            <span className="time-type">швидкий контроль</span>
                        </button>
                        <button
                            className="time-button"
                            onClick={() => handleTimeControlClick(900, 10)}
                        >
                            15+10
                            <span className="time-type">швидкий контроль</span>
                        </button>
                        <button
                            className="time-button"
                            onClick={() => handleTimeControlClick(1500, 10)}
                        >
                            25+10
                            <span className="time-type">швидкий контроль</span>
                        </button>
                        <button
                            className="time-button"
                            onClick={() => handleTimeControlClick(2100, 30)}
                        >
                            35+30
                            <span className="time-type">швидкий контроль</span>
                        </button>
                    </div>
                </div>
            </div>
            <Footer />
        </div>
    );
};

export default HomePage;
