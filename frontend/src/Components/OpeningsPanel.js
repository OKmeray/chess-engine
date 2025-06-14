import React, { useState, useEffect } from "react";

const OpeningsPanel = (setSelectedVariations, selectedVariations) => {
    // State for data, toggles, etc.
    const [openingsData, setOpeningsData] = useState([]);
    const [showOpenings, setShowOpenings] = useState(false);
    const [openOpenings, setOpenOpenings] = useState([]); // which openings are expanded

    useEffect(() => {
        // Same fetch logic you had
        fetch(`${process.env.REACT_APP_BACKEND_URL}games/openings`)
            .then((response) => response.json())
            .then((data) => setOpeningsData(data.openings || []))
            .catch((error) =>
                console.error("Error fetching openings data:", error)
            );
    }, []);

    // Toggles whether an opening is "expanded" (showing its variations)
    const toggleOpening = (openingName) => {
        setOpenOpenings((prev) =>
            prev.includes(openingName)
                ? prev.filter((name) => name !== openingName)
                : [...prev, openingName]
        );
    };

    // Toggles whether a variation is selected
    const toggleVariation = (variationName) => {
        setSelectedVariations((prev) =>
            prev.includes(variationName)
                ? prev.filter((name) => name !== variationName)
                : [...prev, variationName]
        );
    };

    // If an opening is "active" if any of its variations is selected
    const isOpeningActive = (opening) => {
        return opening.variations.some((variation) =>
            selectedVariations.includes(variation.name)
        );
    };

    return (
        <div className="openings-panel">
            <div className="openings-toggle-wrapper">
                <button
                    onClick={() => setShowOpenings(!showOpenings)}
                    className={showOpenings ? "active-opening" : ""}
                >
                    Select Openings
                </button>
                <button
                    onClick={() => setShowOpenings(false)}
                    className={!showOpenings ? "active-opening" : ""}
                >
                    All Openings
                </button>
            </div>

            {showOpenings && (
                <div className="openings-wrapper">
                    <h3>Openings</h3>
                    {openingsData.length > 0 ? (
                        openingsData.map((opening, index) => (
                            <div key={index}>
                                <button
                                    onClick={() => toggleOpening(opening.name)}
                                    className={
                                        isOpeningActive(opening)
                                            ? "active-opening"
                                            : ""
                                    }
                                >
                                    {opening.name}
                                </button>

                                {/* If this opening is expanded, show its variations */}
                                {openOpenings.includes(opening.name) && (
                                    <div className="variations-wrapper">
                                        {opening.variations.map(
                                            (variation, vIndex) => (
                                                <button
                                                    key={vIndex}
                                                    className={
                                                        selectedVariations.includes(
                                                            variation.name
                                                        )
                                                            ? "active-variation"
                                                            : ""
                                                    }
                                                    onClick={() =>
                                                        toggleVariation(
                                                            variation.name
                                                        )
                                                    }
                                                >
                                                    {variation.name}
                                                </button>
                                            )
                                        )}
                                    </div>
                                )}
                            </div>
                        ))
                    ) : (
                        <p>Loading openings...</p>
                    )}
                </div>
            )}
        </div>
    );
};

export default OpeningsPanel;
