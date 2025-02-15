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
                <div className="hero">
                    <h1>Welcome Back to Our Chess Engine</h1>
                    <p>
                        Challenge yourself, play against friends, or test your
                        skills against our advanced AI.
                    </p>
                </div>
                <div className="time-controls">
                    <h2>Select Time Control</h2>
                    <div className="time-buttons">
                        <button
                            className="time-button"
                            onClick={() => handleTimeControlClick(180, 2)}
                        >
                            3+2
                            <span className="time-type">blitz</span>
                        </button>
                        <button
                            className="time-button"
                            onClick={() => handleTimeControlClick(300, 0)}
                        >
                            5+0
                            <span className="time-type">blitz</span>
                        </button>
                        <button
                            className="time-button"
                            onClick={() => handleTimeControlClick(300, 3)}
                        >
                            5+3
                            <span className="time-type">blitz</span>
                        </button>
                        <button
                            className="time-button"
                            onClick={() => handleTimeControlClick(600, 0)}
                        >
                            10+0
                            <span className="time-type">rapid</span>
                        </button>
                    </div>
                </div>
                {/* <div className="recent-games">
          <h2>Recent Games</h2>
          <div className="game-list">
            <div className="game-item">
              <h3>Player1 vs Player2</h3>
              <p>Result: 1-0</p>
            </div>
            <div className="game-item">
              <h3>Player3 vs Player4</h3>
              <p>Result: 0.5-0.5</p>
            </div>
            <div className="game-item">
              <h3>Player5 vs Player6</h3>
              <p>Result: 0-1</p>
            </div>
          </div>
        </div> */}
                <div className="testimonials">
                    <h2>What Our Users Say</h2>
                    <div className="testimonial-list">
                        <div className="testimonial-item">
                            <p>
                                "This platform has transformed my chess game.
                                The AI is challenging and the features are
                                top-notch!"
                            </p>
                            <h4>- Alex Johnson</h4>
                        </div>
                        <div className="testimonial-item">
                            <p>
                                "I love playing against my friends online. The
                                multiplayer mode is smooth and fun."
                            </p>
                            <h4>- Maria Gonzalez</h4>
                        </div>
                    </div>
                </div>
            </div>
            <Footer />
        </div>
    );
};

export default HomePage;
