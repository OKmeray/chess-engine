import React, { useState } from 'react';
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';

import blackKing from '../Assets/Images/king_black.svg';
import blackQueen from '../Assets/Images/queen_black.svg';
import blackRook from '../Assets/Images/rook_black.svg';
import blackBishop from '../Assets/Images/bishop_black.svg';
import blackKnight from '../Assets/Images/knight_black.svg';
import blackPawn from '../Assets/Images/pawn_black.svg';
import whiteKing from '../Assets/Images/king_white.svg';
import whiteQueen from '../Assets/Images/queen_white.svg';
import whiteRook from '../Assets/Images/rook_white.svg';
import whiteBishop from '../Assets/Images/bishop_white.svg';
import whiteKnight from '../Assets/Images/knight_white.svg';
import whitePawn from '../Assets/Images/pawn_white.svg';

const ChessBoard = ({fen}) => {

    const squareSize = 80;  // px

    const boardWrapperStyles = {
        width: squareSize * 8 + "px",
        height: squareSize * 8 + "px",
        position: "relative",
        align: "center"
    }

    const pieceStyles = {
        position: 'absolute',
        top: 0,
        left: 0,
        zIndex: 10000,
        width: squareSize + "px",
    }

    //fen = "3k4/8/5q2/8/2Q5/4P3/8/3K4 w - - 0 1";
    //     /* get position of the pieces and set it in array */
    const piecePlacement = fen.split(' ')[0];

    function isDigit(char) {
        if (char >= '0' && char <= '9') {
            return true;
        }
        return false;
    }

    function getImageFromFenSymbol(char) {
        switch (char) {
            case 'k':
                return blackKing;
            case 'q':
                return blackQueen;
            case 'r':
                return blackRook;
            case 'b':
                return blackBishop;
            case 'n':
                return blackKnight;
            case 'p':
                return blackPawn;
            case 'K':
                return whiteKing;
            case 'Q':
                return whiteQueen;
            case 'R':
                return whiteRook;
            case 'B':
                return whiteBishop;
            case 'N':
                return whiteKnight;
            case 'P':
                return whitePawn
        }
    }

    let pieces = [];
    let squareIndex = 0;
    for (let i = 0; i < piecePlacement.length; i++) {

        if (isDigit(piecePlacement[i]))
        {
            squareIndex += parseInt(piecePlacement[i]);
        }
        else if (piecePlacement[i] === "/") {
            continue;
        }
        else {
            pieces.push({
                "piece": piecePlacement[i],
                "offset": squareIndex
            });
            squareIndex += 1;
        }
        
    }

    console.log(pieces);

    return (
        <div className='board-image-wrapper' style={boardWrapperStyles}>
            <svg width={squareSize * 8} height={squareSize * 8} xmlns="http://www.w3.org/2000/svg">
                {/* Background rectangle */}
                <rect width="100%" height="100%" fill="#e0e0e0" /> {/*fill="#e0e0e0"*/}

                {/* Squares */}
                <g id="squares">
                    {Array.from({ length: 8 }, (_, row) => (
                        Array.from({ length: 8 }, (_, col) => (
                            <rect
                                x={col * squareSize}
                                y={row * squareSize}
                                width={squareSize}
                                height={squareSize}
                                fill={(row + col) % 2 === 0 ? '#f0d9b5' : '#b58863'}
                                key={`${row}-${col}`}
                            />
                        ))
                    ))}
                </g>
            </svg>

            {/* pieces */}
            {pieces.map((pieceObj, index) => (
                <img
                    key={index}
                    src={getImageFromFenSymbol(pieceObj.piece)}
                    style={{
                        position: 'absolute',
                        top: Math.floor(pieceObj.offset / 8) * squareSize,
                        left: (pieceObj.offset % 8) * squareSize,
                        width: squareSize
                    }}
                    alt={`Chess Piece ${pieceObj.piece}`}
                />
            ))}

        </div>
    );
}


export default ChessBoard;