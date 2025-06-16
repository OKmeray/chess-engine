import React from "react";
import "./GameResultPopup.css";

const GameResultPopup = ({ result, onClose }) => {
    const getMessage = () => {
        switch (result) {
            case "win":
                return "Компʼютер переміг!";
            case "loss":
                return "Ви перемогли!";
            case "draw":
                return "Нічия!";
            case "timeoutWin":
                return "Час комп'ютера вичерпано. Ви виграли!";
            case "timeoutLoss":
                return "Ваш час вичерпано. Ви програли.";
            default:
                return "";
        }
    };

    return (
        <div className="game-result-popup-overlay">
            <div className="game-result-popup">
                <h2>{getMessage()}</h2>
                <button onClick={onClose} className="popup-close-button">
                    Нова гра
                </button>
            </div>
        </div>
    );
};

export default GameResultPopup;
