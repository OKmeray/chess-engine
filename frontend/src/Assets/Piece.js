import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { movePiece } from '../Store/actions'; // Import your action creator
import squareSize from '../Assets/variables';
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

const Piece = ({ piece, square }) => {
  const dispatch = useDispatch();
  const isUserTurn = useSelector(state => state.isUserTurn); // Assuming structure


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
    
      function getImageStyles(piece, square) {
        return {
          position: "absolute",
          top: `${Math.floor(square / 8) * squareSize}px`,
          left: `${(square % 8) * squareSize}px`
        }
      }

  const onDragStart = e => {
    if (!isUserTurn) return; // Prevent drag if not user's turn
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/plain', `${piece},${square}`);
    setTimeout(() => {
      e.target.style.display = "none";
    }, 0);
  };

  const onDragEnd = e => {
    e.target.style.display = "";
  };

  const styles = {
    position: "absolute",
    top: `${Math.floor(square / 8) * squareSize}px`,
    left: `${(square % 8) * squareSize}px`
  };

  return (
    <img 
      src={getImageFromFenSymbol(piece)}
      width={squareSize}
      height={squareSize}
      style={styles}
      draggable={isUserTurn} // Only draggable if it's the user's turn
      onDragStart={onDragStart}
      onDragEnd={onDragEnd}
    />
  );
};

export default Piece;
