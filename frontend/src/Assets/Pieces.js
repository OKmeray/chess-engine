import React, { useState, useRef } from 'react';
import Piece from '../Assets/Piece';
import squareSize from '../Assets/variables';

const Pieces = ({ pieces, onPieceDrop }) => {
    const ref = useRef();

    const piecesWrapperStyles = {
        width: squareSize * 8 + "px",
        height: squareSize * 8 + "px",
        position: "absolute",
    };

    // Calculate the board square based on mouse position
    const calculateSquare = (clientX, clientY) => {
        const { left, top } = ref.current.getBoundingClientRect();
        const file = Math.floor((clientX - left) / squareSize);
        const rank = Math.floor((clientY - top) / squareSize); // Adjust if your board is flipped
        return rank * 8 + file;
    };

    // Handle the drop event when a piece is dropped
    const handleDrop = (event) => {
        event.preventDefault();
        const data = event.dataTransfer.getData('text/plain');
        const [piece, fromSquare] = data.split(',');
        const toSquare = calculateSquare(event.clientX, event.clientY);

        // Update pieces array and notify parent component
        onPieceDrop(fromSquare, toSquare);
    };

    // Allow drag over to enable drop
    const onDragOver = (event) => {
        event.preventDefault();
    };

    return (
        <div ref={ref} style={piecesWrapperStyles} onDrop={handleDrop} onDragOver={onDragOver}>
            {pieces.map((piece, index) => (
                <Piece key={index} piece={piece.piece} square={piece.offset} />
            ))}
        </div>
    );
};

export default Pieces;
