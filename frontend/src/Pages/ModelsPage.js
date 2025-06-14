import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import Header from "../components/Header";
import Footer from "../components/Footer";
import "./ModelsPage.css";
import Modal from "react-modal";

// Hardcoded data for now (will be fetched from backend later)
const modelsData = [
    {
        HEXID: "0803",
        ID: 1,
        Filters: 64,
        ResBlocks: 6,
        BatchSize: 128,
        LabelSmoothing: 0,
        LR: 0.001,
        SEBlocks: true,
        LossWeights: "policy: 0.5; value: 0.5",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 1866000,
        PuzzlePositions: 375000,
        test_policy: 2.5954,
        test_value: 0.3213,
        mae: 0.4173,
        test_top1: 0.3059,
        test_top3: 0.5296,
        test_top5: 0.6474,
    },
    {
        HEXID: "0f39",
        ID: 2,
        Filters: 64,
        ResBlocks: 6,
        BatchSize: 128,
        LabelSmoothing: 0.002,
        LR: 0.001,
        SEBlocks: true,
        LossWeights: "policy: 0.5; value: 0.5",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 1866000,
        PuzzlePositions: 375000,
        test_policy: 2.626,
        test_value: 0.3179,
        mae: 0.4249,
        test_top1: 0.2989,
        test_top3: 0.5229,
        test_top5: 0.6407,
    },
    {
        HEXID: "110f",
        ID: 3,
        Filters: 64,
        ResBlocks: 6,
        BatchSize: 128,
        LabelSmoothing: 0.005,
        LR: 0.001,
        SEBlocks: true,
        LossWeights: "policy: 0.5; value: 0.5",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 1866000,
        PuzzlePositions: 375000,
        test_policy: 2.6368,
        test_value: 0.337,
        mae: 0.4587,
        test_top1: 0.3129,
        test_top3: 0.5429,
        test_top5: 0.6619,
    },
    {
        HEXID: "1a45",
        ID: 4,
        Filters: 64,
        ResBlocks: 6,
        BatchSize: 128,
        LabelSmoothing: 0,
        LR: 0.001,
        SEBlocks: true,
        LossWeights: "policy: 0.7; value: 0.3",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 1866000,
        PuzzlePositions: 375000,
        test_policy: 2.5603,
        test_value: 0.3279,
        mae: 0.4407,
        test_top1: 0.3101,
        test_top3: 0.5429,
        test_top5: 0.6629,
    },
    {
        HEXID: "21fb",
        ID: 5,
        Filters: 64,
        ResBlocks: 6,
        BatchSize: 128,
        LabelSmoothing: 0,
        LR: 0.001,
        SEBlocks: true,
        LossWeights: "policy: 0.8; value: 0.2",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 1866000,
        PuzzlePositions: 375000,
        test_policy: 2.5751,
        test_value: 0.3826,
        mae: 0.4919,
        test_top1: 0.309,
        test_top3: 0.5433,
        test_top5: 0.6633,
    },
    {
        HEXID: "23a6",
        ID: 6,
        Filters: 64,
        ResBlocks: 6,
        BatchSize: 128,
        LabelSmoothing: 0,
        LR: 0.001,
        SEBlocks: true,
        LossWeights: "policy: 0.6; value: 0.4",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 1866000,
        PuzzlePositions: 375000,
        test_policy: 2.6528,
        test_value: 0.3761,
        mae: 0.4709,
        test_top1: 0.3063,
        test_top3: 0.5367,
        test_top5: 0.6543,
    },
    {
        HEXID: "25ab",
        ID: 7,
        Filters: 64,
        ResBlocks: 6,
        BatchSize: 128,
        LabelSmoothing: 0.001,
        LR: 0.001,
        SEBlocks: true,
        LossWeights: "policy: 0.5; value: 0.5",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 1866000,
        PuzzlePositions: 375000,
        test_policy: 2.6113,
        test_value: 0.329,
        mae: 0.4212,
        test_top1: 0.3004,
        test_top3: 0.529,
        test_top5: 0.6497,
    },
    {
        HEXID: "2b52",
        ID: 8,
        Filters: 64,
        ResBlocks: 6,
        BatchSize: 128,
        LabelSmoothing: 0,
        LR: 0.001,
        SEBlocks: true,
        LossWeights: "policy: 0.7; value: 0.3",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 2000000,
        PuzzlePositions: 2000000,
        test_policy: 1.5153,
        test_value: 0.2252,
        mae: 0.2981,
        test_top1: 0.5959,
        test_top3: 0.7617,
        test_top5: 0.8267,
    },
    {
        HEXID: "2c52",
        ID: 9,
        Filters: 64,
        ResBlocks: 6,
        BatchSize: 128,
        LabelSmoothing: 0,
        LR: 0.001,
        SEBlocks: true,
        LossWeights: "policy: 0.7; value: 0.3",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 1000000,
        PuzzlePositions: 3000000,
        test_policy: 1.2248,
        test_value: 0.1639,
        mae: 0.2469,
        test_top1: 0.6888,
        test_top3: 0.8174,
        test_top5: 0.8618,
    },
    {
        HEXID: "2c6f",
        ID: 10,
        Filters: 64,
        ResBlocks: 6,
        BatchSize: 128,
        LabelSmoothing: 0,
        LR: 0.001,
        SEBlocks: true,
        LossWeights: "policy: 0.7; value: 0.3",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 3000000,
        PuzzlePositions: 1000000,
        test_policy: 2.1077,
        test_value: 0.247,
        mae: 0.3478,
        test_top1: 0.4331,
        test_top3: 0.6425,
        test_top5: 0.7419,
    },
    {
        HEXID: "30fe",
        ID: 11,
        Filters: 64,
        ResBlocks: 6,
        BatchSize: 128,
        LabelSmoothing: 0,
        LR: 0.001,
        SEBlocks: true,
        LossWeights: "policy: 0.7; value: 0.3",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 1400000,
        PuzzlePositions: 2600000,
        test_policy: 1.567,
        test_value: 0.2079,
        mae: 0.3074,
        test_top1: 0.589,
        test_top3: 0.75,
        test_top5: 0.8147,
    },
    {
        HEXID: "3171",
        ID: 12,
        Filters: 64,
        ResBlocks: 6,
        BatchSize: 128,
        LabelSmoothing: 0,
        LR: 0.001,
        SEBlocks: true,
        LossWeights: "policy: 0.7; value: 0.3",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 2600000,
        PuzzlePositions: 1400000,
        test_policy: 1.8207,
        test_value: 0.2161,
        mae: 0.3112,
        test_top1: 0.5258,
        test_top3: 0.6943,
        test_top5: 0.7721,
    },
    {
        HEXID: "4380",
        ID: 14,
        Filters: 64,
        ResBlocks: 6,
        BatchSize: 128,
        LabelSmoothing: 0,
        LR: 0.001,
        SEBlocks: false,
        LossWeights: "policy: 0.7; value: 0.3",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 1400000,
        PuzzlePositions: 2600000,
        test_policy: 1.4548,
        test_value: 0.1668,
        mae: 0.2697,
        test_top1: 0.6057,
        test_top3: 0.7668,
        test_top5: 0.8369,
    },
    {
        HEXID: "4389",
        ID: 15,
        Filters: 64,
        ResBlocks: 6,
        BatchSize: 128,
        LabelSmoothing: 0,
        LR: 0.001,
        SEBlocks: false,
        LossWeights: "policy: 0.7; value: 0.3",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 2000000,
        PuzzlePositions: 2000000,
        test_policy: 1.6846,
        test_value: 0.2136,
        mae: 0.291,
        test_top1: 0.5443,
        test_top3: 0.7206,
        test_top5: 0.8014,
    },
    {
        HEXID: "4d51",
        ID: 16,
        Filters: 64,
        ResBlocks: 6,
        BatchSize: 128,
        LabelSmoothing: 0,
        LR: 0.001,
        SEBlocks: false,
        LossWeights: "policy: 0.7; value: 0.3",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 2600000,
        PuzzlePositions: 1400000,
        test_policy: 1.7399,
        test_value: 0.19,
        mae: 0.2779,
        test_top1: 0.5216,
        test_top3: 0.7106,
        test_top5: 0.7983,
    },
    {
        HEXID: "5d7a",
        ID: 17,
        Filters: 64,
        ResBlocks: 4,
        BatchSize: 128,
        LabelSmoothing: 0,
        LR: 0.001,
        SEBlocks: true,
        LossWeights: "policy: 0.7; value: 0.3",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 2000000,
        PuzzlePositions: 2000000,
        test_policy: 1.5923,
        test_value: 0.2299,
        mae: 0.269,
        test_top1: 0.58,
        test_top3: 0.7376,
        test_top5: 0.8069,
    },
    {
        HEXID: "5da8",
        ID: 18,
        Filters: 32,
        ResBlocks: 4,
        BatchSize: 128,
        LabelSmoothing: 0,
        LR: 0.001,
        SEBlocks: true,
        LossWeights: "policy: 0.7; value: 0.3",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 2000000,
        PuzzlePositions: 2000000,
        test_policy: 1.6877,
        test_value: 0.2344,
        mae: 0.3234,
        test_top1: 0.5535,
        test_top3: 0.725,
        test_top5: 0.7964,
    },
    {
        HEXID: "78fc",
        ID: 19,
        Filters: 32,
        ResBlocks: 2,
        BatchSize: 128,
        LabelSmoothing: 0,
        LR: 0.001,
        SEBlocks: true,
        LossWeights: "policy: 0.7; value: 0.3",
        DropoutRate: 0,
        L2: 0.0005,
        Epochs: 5,
        SimplePositions: 2000000,
        PuzzlePositions: 2000000,
        test_policy: 1.7923,
        test_value: 0.248,
        mae: 0.3757,
        test_top1: 0.5267,
        test_top3: 0.7055,
        test_top5: 0.7791,
    },
];

const ModelsPage = () => {
    const [selectedModel, setSelectedModel] = useState(null);
    const [modalIsOpen, setModalIsOpen] = useState(false);
    const [filteredModels, setFilteredModels] = useState(modelsData);
    const [sortConfig, setSortConfig] = useState({
        key: null,
        direction: "asc",
    });
    const [filters, setFilters] = useState({
        minTop1: 0,
        maxMAE: 1.0,
        hasSEBlocks: false,
    });

    // whenever filters or sortConfig changes
    useEffect(() => {
        let result = [...modelsData];

        // Apply filters
        result = result.filter((model) => {
            return (
                model.test_top1 * 100 >= filters.minTop1 &&
                model.mae <= filters.maxMAE &&
                (!filters.hasSEBlocks || model.SEBlocks)
            );
        });

        // Apply sorting
        if (sortConfig.key) {
            result.sort((a, b) => {
                if (a[sortConfig.key] < b[sortConfig.key]) {
                    return sortConfig.direction === "asc" ? -1 : 1;
                }
                if (a[sortConfig.key] > b[sortConfig.key]) {
                    return sortConfig.direction === "asc" ? 1 : -1;
                }
                return 0;
            });
        }

        setFilteredModels(result);
    }, [filters, sortConfig]);

    useEffect(() => {
        if (modalIsOpen) {
            document.body.classList.add("body-no-scroll");
        } else {
            document.body.classList.remove("body-no-scroll");
        }

        return () => {
            document.body.classList.remove("body-no-scroll");
        };
    }, [modalIsOpen]);

    const openModal = (model) => {
        setSelectedModel(model);
        setModalIsOpen(true);
    };

    const closeModal = () => {
        setModalIsOpen(false);
    };

    const formatNumber = (num) => {
        if (typeof num === "number") {
            return num.toLocaleString();
        }
        return num;
    };

    const handleSort = (key) => {
        let direction = "asc";
        if (sortConfig.key === key && sortConfig.direction === "asc") {
            direction = "desc";
        }
        setSortConfig({ key, direction });
    };

    const handleFilterChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFilters((prev) => ({
            ...prev,
            [name]: type === "checkbox" ? checked : parseFloat(value),
        }));
    };

    const resetFilters = () => {
        setFilters({
            minTop1: 0,
            maxMAE: 1.0,
            hasSEBlocks: false,
        });
        setSortConfig({ key: null, direction: "asc" });
    };

    return (
        <div className="page-container">
            <Header />
            <div className="content-container">
                <div className="white-wrapper">
                    <h1>Доступні шахові моделі</h1>
                    <p>Натисніть на модель, щоб побачити деталі моделі</p>
                </div>

                <div className="white-wrapper">
                    <div className="filter-group">
                        <label>
                            Мін. точність Top-1 (%):
                            <input
                                type="number"
                                name="minTop1"
                                min="0"
                                max="100"
                                step="0.1"
                                value={filters.minTop1}
                                onChange={handleFilterChange}
                            />
                        </label>
                    </div>

                    <div className="filter-group">
                        <label>
                            Макс. значення MAE:
                            <input
                                type="number"
                                name="maxMAE"
                                min="0"
                                max="1"
                                step="0.01"
                                value={filters.maxMAE}
                                onChange={handleFilterChange}
                            />
                        </label>
                    </div>

                    <div className="filter-group">
                        <label>
                            <input
                                type="checkbox"
                                name="hasSEBlocks"
                                checked={filters.hasSEBlocks}
                                onChange={handleFilterChange}
                            />
                            Тільки моделі з SE-блоками
                        </label>
                    </div>

                    <div className="sort-buttons">
                        <span>Сортувати за:</span>
                        <button onClick={() => handleSort("test_top1")}>
                            Точністю Top-1{" "}
                            {sortConfig.key === "test_top1" &&
                                `(${sortConfig.direction})`}
                        </button>
                        <button onClick={() => handleSort("mae")}>
                            Значенням MAE{" "}
                            {sortConfig.key === "mae" &&
                                `(${sortConfig.direction})`}
                        </button>
                        <button onClick={() => handleSort("ID")}>
                            ID{" "}
                            {sortConfig.key === "ID" &&
                                `(${sortConfig.direction})`}
                        </button>
                        <button onClick={resetFilters}>Скинути фільтри</button>
                    </div>
                </div>

                <div className="models-grid">
                    {filteredModels.length > 0 ? (
                        filteredModels.map((model) => (
                            <div
                                key={model.HEXID}
                                className="model-card"
                                onClick={() => openModal(model)}
                            >
                                <div className="model-card-header">
                                    <h3>Модель #{model.ID}</h3>
                                    <span className="model-hexid">
                                        {model.HEXID}
                                    </span>
                                </div>
                                <div className="model-card-body">
                                    <div className="model-stats">
                                        <div className="stat-item">
                                            <span className="stat-label">
                                                Точність Top-1
                                            </span>
                                            <span className="stat-value">
                                                {(
                                                    model.test_top1 * 100
                                                ).toFixed(1)}
                                                %
                                            </span>
                                        </div>
                                        <div className="stat-item">
                                            <span className="stat-label">
                                                Точність Top-5
                                            </span>
                                            <span className="stat-value">
                                                {(
                                                    model.test_top5 * 100
                                                ).toFixed(1)}
                                                %
                                            </span>
                                        </div>
                                        <div className="stat-item">
                                            <span className="stat-label">
                                                Значення MAE
                                            </span>
                                            <span className="stat-value">
                                                {model.mae.toFixed(4)}
                                            </span>
                                        </div>
                                    </div>
                                </div>
                                <div className="model-card-footer">
                                    <span className="model-config">
                                        Фільтри: {model.Filters} | Залишкові
                                        блоки: {model.ResBlocks}
                                    </span>
                                    <a
                                        href={`/game?model=${model.HEXID}`}
                                        className="play-button"
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            // Reset scroll position before navigation
                                            window.scrollTo(0, 0);
                                            // Optional: Add smooth scroll
                                            window.scrollTo({
                                                top: 0,
                                                behavior: "smooth",
                                            });
                                        }}
                                    >
                                        Грати з моделлю
                                    </a>
                                </div>
                            </div>
                        ))
                    ) : (
                        <div className="no-results">
                            За заданими фільтрами моделей не знайдено. Спробуйте
                            інші фільтри.
                        </div>
                    )}
                </div>
            </div>

            <Modal
                isOpen={modalIsOpen}
                onRequestClose={closeModal}
                contentLabel="Model Details"
                className="model-modal"
                overlayClassName="model-modal-overlay"
                shouldCloseOnOverlayClick={true}
                onAfterOpen={() => {
                    // Add event listener to prevent wheel events from bubbling
                    const modal = document.querySelector(".model-modal");
                    if (modal) {
                        modal.addEventListener("wheel", (e) => {
                            e.stopPropagation();
                        });
                    }
                }}
                onAfterClose={() => {
                    // Clean up event listener
                    const modal = document.querySelector(".model-modal");
                    if (modal) {
                        modal.removeEventListener("wheel", (e) => {
                            e.stopPropagation();
                        });
                    }
                }}
            >
                {selectedModel && (
                    <div className="modal-content">
                        <div className="modal-header">
                            <h2>Модель #{selectedModel.ID}</h2>
                            <button
                                onClick={closeModal}
                                className="close-button"
                            >
                                &times;
                            </button>
                        </div>

                        <div className="modal-body">
                            <div className="model-details-grid">
                                <div className="details-section">
                                    <h3>Налаштування моделі</h3>
                                    <ul>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <div className="info-icon">
                                                    <svg
                                                        width="14"
                                                        height="14"
                                                        viewBox="0 0 18 18"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                            fill="black"
                                                        />
                                                        <circle
                                                            cx="9"
                                                            cy="9"
                                                            r="8.5"
                                                            stroke="black"
                                                        />
                                                    </svg>
                                                    <span className="tooltip-text">
                                                        Унікальний ідентифікатор
                                                        моделі у
                                                        шістнадцятковому
                                                        форматі.
                                                        Використовується для
                                                        ідентифікації моделі в
                                                        системі.
                                                    </span>
                                                </div>
                                                <strong>ID:</strong>
                                            </div>
                                            {selectedModel.HEXID}
                                        </li>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <div className="info-icon">
                                                    <svg
                                                        width="14"
                                                        height="14"
                                                        viewBox="0 0 18 18"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                            fill="black"
                                                        />
                                                        <circle
                                                            cx="9"
                                                            cy="9"
                                                            r="8.5"
                                                            stroke="black"
                                                        />
                                                    </svg>
                                                    <span className="tooltip-text">
                                                        Кількість фільтрів у
                                                        згорткових шарах. Більше
                                                        фільтрів може покращити
                                                        якість моделі, але
                                                        збільшує вимоги до
                                                        обчислювальних ресурсів.
                                                    </span>
                                                </div>
                                                <strong>Фільтри:</strong>
                                            </div>
                                            {selectedModel.Filters}
                                        </li>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <div className="info-icon">
                                                    <svg
                                                        width="14"
                                                        height="14"
                                                        viewBox="0 0 18 18"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                            fill="black"
                                                        />
                                                        <circle
                                                            cx="9"
                                                            cy="9"
                                                            r="8.5"
                                                            stroke="black"
                                                        />
                                                    </svg>
                                                    <span className="tooltip-text">
                                                        Кількість залишкових
                                                        блоків у архітектурі. Ці
                                                        блоки допомагають
                                                        навчати глибші моделі,
                                                        запобігаючи проблемі
                                                        зникаючого градієнта.
                                                    </span>
                                                </div>
                                                <strong>
                                                    Залишкові блоки:
                                                </strong>
                                            </div>
                                            {selectedModel.ResBlocks}
                                        </li>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <div className="info-icon">
                                                    <svg
                                                        width="14"
                                                        height="14"
                                                        viewBox="0 0 18 18"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                            fill="black"
                                                        />
                                                        <circle
                                                            cx="9"
                                                            cy="9"
                                                            r="8.5"
                                                            stroke="black"
                                                        />
                                                    </svg>
                                                    <span className="tooltip-text">
                                                        Розмір пакету даних для
                                                        одного кроку навчання.
                                                        Великі пакети можуть
                                                        пришвидшити навчання,
                                                        але вимагають більше
                                                        пам'яті GPU.
                                                    </span>
                                                </div>
                                                <strong>Розмір пакетів:</strong>
                                            </div>
                                            {selectedModel.BatchSize}
                                        </li>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <div className="info-icon">
                                                    <svg
                                                        width="14"
                                                        height="14"
                                                        viewBox="0 0 18 18"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                            fill="black"
                                                        />
                                                        <circle
                                                            cx="9"
                                                            cy="9"
                                                            r="8.5"
                                                            stroke="black"
                                                        />
                                                    </svg>
                                                    <span className="tooltip-text">
                                                        Швидкість навчання
                                                        нейромережі. Визначає
                                                        розмір кроку, який
                                                        робить оптимізатор при
                                                        оновленні вагів моделі.
                                                    </span>
                                                </div>
                                                <strong>
                                                    Коефіцієнт навчання:
                                                </strong>
                                            </div>
                                            {selectedModel.LR}
                                        </li>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <div className="info-icon">
                                                    <svg
                                                        width="14"
                                                        height="14"
                                                        viewBox="0 0 18 18"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                            fill="black"
                                                        />
                                                        <circle
                                                            cx="9"
                                                            cy="9"
                                                            r="8.5"
                                                            stroke="black"
                                                        />
                                                    </svg>
                                                    <span className="tooltip-text">
                                                        Чи використовуються
                                                        SE-блоки
                                                        (Squeeze-and-Excitation).
                                                        Ці блоки покращують
                                                        якість моделі,
                                                        дозволяючи їй динамічно
                                                        перерозподіляти увагу
                                                        між різними каналами.
                                                    </span>
                                                </div>
                                                <strong>
                                                    Наявність SE-блоків:
                                                </strong>
                                            </div>
                                            {selectedModel.SEBlocks
                                                ? "Так"
                                                : "Ні"}
                                        </li>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <div className="info-icon">
                                                    <svg
                                                        width="14"
                                                        height="14"
                                                        viewBox="0 0 18 18"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                            fill="black"
                                                        />
                                                        <circle
                                                            cx="9"
                                                            cy="9"
                                                            r="8.5"
                                                            stroke="black"
                                                        />
                                                    </svg>
                                                    <span className="tooltip-text">
                                                        Ймовірність вимкнення
                                                        нейронів під час
                                                        навчання. Допомагає
                                                        запобігти перенавчанню,
                                                        змушуючи модель бути
                                                        більш стійкою.
                                                    </span>
                                                </div>
                                                <strong>Рівень dropout:</strong>
                                            </div>
                                            {selectedModel.DropoutRate}
                                        </li>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <div className="info-icon">
                                                    <svg
                                                        width="14"
                                                        height="14"
                                                        viewBox="0 0 18 18"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                            fill="black"
                                                        />
                                                        <circle
                                                            cx="9"
                                                            cy="9"
                                                            r="8.5"
                                                            stroke="black"
                                                        />
                                                    </svg>
                                                    <span className="tooltip-text">
                                                        Сила L2 регуляризації.
                                                        Додає штраф за великі
                                                        значення вагів,
                                                        запобігаючи перенавчанню
                                                        моделі.
                                                    </span>
                                                </div>
                                                <strong>
                                                    L2 регуляризація:
                                                </strong>
                                            </div>
                                            {selectedModel.L2}
                                        </li>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <div className="info-icon">
                                                    <svg
                                                        width="14"
                                                        height="14"
                                                        viewBox="0 0 18 18"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                            fill="black"
                                                        />
                                                        <circle
                                                            cx="9"
                                                            cy="9"
                                                            r="8.5"
                                                            stroke="black"
                                                        />
                                                    </svg>
                                                    <span className="tooltip-text">
                                                        Кількість повних
                                                        проходів через
                                                        навчальний набір. Більше
                                                        епох зазвичай покращує
                                                        якість моделі, але
                                                        збільшує час навчання.
                                                    </span>
                                                </div>
                                                <strong>Кількість епох:</strong>
                                            </div>
                                            {selectedModel.Epochs}
                                        </li>
                                    </ul>
                                </div>
                                <div className="details-section">
                                    <h3>Показники ефективності</h3>
                                    <ul>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <div className="info-icon">
                                                    <svg
                                                        width="14"
                                                        height="14"
                                                        viewBox="0 0 18 18"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                            fill="black"
                                                        />
                                                        <circle
                                                            cx="9"
                                                            cy="9"
                                                            r="8.5"
                                                            stroke="black"
                                                        />
                                                    </svg>
                                                    <span className="tooltip-text">
                                                        Втрати для передбачення
                                                        ймовірностей ходів
                                                        (кросс-ентропія). Чим
                                                        менше значення, тим
                                                        краще модель передбачає
                                                        правильні ходи.
                                                    </span>
                                                </div>
                                                <strong>Втрати Policy:</strong>
                                            </div>
                                            {selectedModel.test_policy.toFixed(
                                                4
                                            )}
                                        </li>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <div className="info-icon">
                                                    <svg
                                                        width="14"
                                                        height="14"
                                                        viewBox="0 0 18 18"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                            fill="black"
                                                        />
                                                        <circle
                                                            cx="9"
                                                            cy="9"
                                                            r="8.5"
                                                            stroke="black"
                                                        />
                                                    </svg>
                                                    <span className="tooltip-text">
                                                        Втрати для оцінки
                                                        позиції
                                                        (середньоквадратична
                                                        помилка). Відображає,
                                                        наскільки точно модель
                                                        оцінює, хто виграє в
                                                        поточній позиції.
                                                    </span>
                                                </div>
                                                <strong>Втрати Value:</strong>
                                            </div>
                                            {selectedModel.test_value.toFixed(
                                                4
                                            )}
                                        </li>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <div className="info-icon">
                                                    <svg
                                                        width="14"
                                                        height="14"
                                                        viewBox="0 0 18 18"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                            fill="black"
                                                        />
                                                        <circle
                                                            cx="9"
                                                            cy="9"
                                                            r="8.5"
                                                            stroke="black"
                                                        />
                                                    </svg>
                                                    <span className="tooltip-text">
                                                        Середня абсолютна
                                                        похибка оцінки позиції.
                                                        Показує середнє
                                                        відхилення передбачення
                                                        моделі від реального
                                                        результату.
                                                    </span>
                                                </div>
                                                <strong>Значення MAE:</strong>
                                            </div>
                                            {selectedModel.mae.toFixed(4)}
                                        </li>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <div className="info-icon">
                                                    <svg
                                                        width="14"
                                                        height="14"
                                                        viewBox="0 0 18 18"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                            fill="black"
                                                        />
                                                        <circle
                                                            cx="9"
                                                            cy="9"
                                                            r="8.5"
                                                            stroke="black"
                                                        />
                                                    </svg>
                                                    <span className="tooltip-text">
                                                        Імовірність, що модель
                                                        правильно визначила
                                                        найкращий хід серед усіх
                                                        можливих. Важливий
                                                        показник якості
                                                        передбачень.
                                                    </span>
                                                </div>
                                                <strong>Точність Top-1:</strong>
                                            </div>
                                            {(
                                                selectedModel.test_top1 * 100
                                            ).toFixed(1)}
                                            %
                                        </li>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <div className="info-icon">
                                                    <svg
                                                        width="14"
                                                        height="14"
                                                        viewBox="0 0 18 18"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                            fill="black"
                                                        />
                                                        <circle
                                                            cx="9"
                                                            cy="9"
                                                            r="8.5"
                                                            stroke="black"
                                                        />
                                                    </svg>
                                                    <span className="tooltip-text">
                                                        Імовірність, що
                                                        правильний хід є серед 3
                                                        найкращих варіантів,
                                                        передбачених моделлю.
                                                        Показує гнучкість
                                                        моделі.
                                                    </span>
                                                </div>
                                                <strong>Точність Top-3:</strong>
                                            </div>
                                            {(
                                                selectedModel.test_top3 * 100
                                            ).toFixed(1)}
                                            %
                                        </li>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <div className="info-icon">
                                                    <svg
                                                        width="14"
                                                        height="14"
                                                        viewBox="0 0 18 18"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                            fill="black"
                                                        />
                                                        <circle
                                                            cx="9"
                                                            cy="9"
                                                            r="8.5"
                                                            stroke="black"
                                                        />
                                                    </svg>
                                                    <span className="tooltip-text">
                                                        Імовірність, що
                                                        правильний хід є серед 5
                                                        найкращих варіантів.
                                                        Високі значення свідчать
                                                        про хорошу якість
                                                        передбачень моделі.
                                                    </span>
                                                </div>
                                                <strong>Точність Top-5:</strong>
                                            </div>
                                            {(
                                                selectedModel.test_top5 * 100
                                            ).toFixed(1)}
                                            %
                                        </li>
                                    </ul>
                                </div>
                                <div className="details-section">
                                    <h3>Тренувальний набір даних</h3>
                                    <ul>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <div className="info-icon">
                                                    <svg
                                                        width="14"
                                                        height="14"
                                                        viewBox="0 0 18 18"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                            fill="black"
                                                        />
                                                        <circle
                                                            cx="9"
                                                            cy="9"
                                                            r="8.5"
                                                            stroke="black"
                                                        />
                                                    </svg>
                                                    <span className="tooltip-text">
                                                        Кількість стандартних
                                                        шахових позицій з партій
                                                        гросмейстерів. Ці
                                                        позиції використовуються
                                                        для навчання моделі
                                                        основним стратегіям гри.
                                                    </span>
                                                </div>
                                                <strong>
                                                    Позиції з партій
                                                    гросмейстерів:
                                                </strong>
                                            </div>
                                            {formatNumber(
                                                selectedModel.SimplePositions
                                            )}
                                        </li>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <div className="info-icon">
                                                    <svg
                                                        width="14"
                                                        height="14"
                                                        viewBox="0 0 18 18"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M8.22727 15V6.27273H9.56818V15H8.22727ZM8.90909 4.81818C8.64773 4.81818 8.42235 4.72917 8.23295 4.55114C8.04735 4.37311 7.95455 4.15909 7.95455 3.90909C7.95455 3.65909 8.04735 3.44508 8.23295 3.26705C8.42235 3.08902 8.64773 3 8.90909 3C9.17045 3 9.39394 3.08902 9.57955 3.26705C9.76894 3.44508 9.86364 3.65909 9.86364 3.90909C9.86364 4.15909 9.76894 4.37311 9.57955 4.55114C9.39394 4.72917 9.17045 4.81818 8.90909 4.81818Z"
                                                            fill="black"
                                                        />
                                                        <circle
                                                            cx="9"
                                                            cy="9"
                                                            r="8.5"
                                                            stroke="black"
                                                        />
                                                    </svg>
                                                    <span className="tooltip-text">
                                                        Кількість тактичних
                                                        позицій з визначеними
                                                        найкращими ходами. Ці
                                                        позиції допомагають
                                                        моделі краще розуміти
                                                        тактичні моменти в грі.
                                                    </span>
                                                </div>
                                                <strong>
                                                    Тактичні позиції:
                                                </strong>
                                            </div>
                                            {formatNumber(
                                                selectedModel.PuzzlePositions
                                            )}
                                        </li>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <strong>
                                                    Загальна кількість:
                                                </strong>
                                            </div>
                                            {formatNumber(
                                                selectedModel.SimplePositions +
                                                    selectedModel.PuzzlePositions
                                            )}
                                        </li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                        <div className="modal-footer">
                            <a
                                href={`/architecture_images/${selectedModel.HEXID}.png`}
                                download={`architecture_${selectedModel.HEXID}.png`}
                                className="base-button hover-selected"
                            >
                                Завантажити архітектуру
                            </a>
                            <button
                                onClick={closeModal}
                                className="modal-close-btn"
                            >
                                Закрити
                            </button>
                        </div>
                    </div>
                )}
            </Modal>

            <Footer />
        </div>
    );
};

export default ModelsPage;
