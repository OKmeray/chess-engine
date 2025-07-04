import { useState, useEffect } from "react";
import Header from "../Components/Header.js";
import Footer from "../Components/Footer.js";
import "./ModelsPage.css";
import Modal from "react-modal";

const ModelsPage = () => {
    const [selectedModel, setSelectedModel] = useState(null);
    const [modalIsOpen, setModalIsOpen] = useState(false);
    const [allModels, setAllModels] = useState([]);
    const [filteredModels, setFilteredModels] = useState([]);
    const [sortConfig, setSortConfig] = useState({
        key: null,
        direction: "за зростанням",
    });
    const [filters, setFilters] = useState({
        minTop1: 0,
        maxMAE: 1.0,
        hasSEBlocks: false,
    });

    useEffect(() => {
        const fetchModels = async () => {
            try {
                const response = await fetch(
                    `${process.env.REACT_APP_BACKEND_URL}api/EngineConfigs`
                );
                if (!response.ok) throw new Error("Failed to fetch");
                const data = await response.json();
                setAllModels(data);
                setFilteredModels(data);
            } catch (error) {
                console.error("Error fetching models:", error);
                // Optionally set some error state here
            }
            // finally {
            //     setIsLoading(false);
            // }
        };

        fetchModels();
    }, []);

    // whenever filters or sortConfig changes
    useEffect(() => {
        if (allModels.length === 0) return;
        let result = [...allModels];

        // Apply filters
        result = result.filter((model) => {
            return (
                model.test_top1 * 100 >= filters.minTop1 &&
                model.mae <= filters.maxMAE &&
                (!filters.hasSEBlocks || model.seBlocks)
            );
        });

        // Apply sorting
        if (sortConfig.key) {
            result.sort((a, b) => {
                if (a[sortConfig.key] < b[sortConfig.key]) {
                    return sortConfig.direction === "за зростанням" ? -1 : 1;
                }
                if (a[sortConfig.key] > b[sortConfig.key]) {
                    return sortConfig.direction === "за зростанням" ? 1 : -1;
                }
                return 0;
            });
        }

        setFilteredModels(result);
    }, [filters, sortConfig, allModels]);

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
        const backendKeyMap = {
            test_top1: "test_top1",
            mae: "mae",
            ID: "configId",
        };

        const backendKey = backendKeyMap[key] || key;

        let direction = "за зростанням";
        if (
            sortConfig.key === key &&
            sortConfig.direction === "за зростанням"
        ) {
            direction = "за спаданням";
        }

        setSortConfig({ key: backendKey, direction });
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
        setSortConfig({ key: null, direction: "за зростанням" });
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
                        <button onClick={resetFilters}>Скинути фільтри</button>
                    </div>
                </div>

                <div className="models-grid">
                    {filteredModels.length > 0 ? (
                        filteredModels.map((model) => (
                            <div
                                key={model.hexid}
                                className="model-card"
                                onClick={() => openModal(model)}
                            >
                                <div className="model-card-header">
                                    <h3>Модель #{model.configId}</h3>
                                    <span className="model-hexid">
                                        {model.hexid}
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
                                        Фільтри: {model.filters} | Залишкові
                                        блоки: {model.resBlocks}
                                    </span>
                                    <a
                                        href={`/game?model=${model.hexid}`}
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
                            <h2>Модель #{selectedModel.configId}</h2>
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
                                                <strong>Ідентифікатор:</strong>
                                            </div>
                                            {selectedModel.hexid}
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
                                            {selectedModel.filters}
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
                                            {selectedModel.resBlocks}
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
                                            {selectedModel.batchSize}
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
                                            {selectedModel.lr}
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
                                            {selectedModel.seBlocks
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
                                            {selectedModel.dropoutRate}
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
                                            {selectedModel.l2}
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
                                            {selectedModel.epochs}
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
                                                selectedModel.simplePositions
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
                                                selectedModel.puzzlePositions
                                            )}
                                        </li>
                                        <li>
                                            <div className="parameter-name-wrapper">
                                                <strong>
                                                    Загальна кількість:
                                                </strong>
                                            </div>
                                            {formatNumber(
                                                selectedModel.simplePositions +
                                                    selectedModel.puzzlePositions
                                            )}
                                        </li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                        <div className="modal-footer">
                            <a
                                href={`/architecture_images/${selectedModel.hexid}.png`}
                                download={`architecture_${selectedModel.hexid}.png`}
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
