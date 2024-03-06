import squareSize from '../Assets/variables';
import Piece from '../Assets/Piece';
import {useState, useRef} from 'react';


const Pieces = (initialPieces) => {

    const ref = useRef();

    const [pieces, setPieces] = useState(initialPieces);
    
    const piecesWrapperStyles = {
        width: squareSize * 8 + "px",
        height: squareSize * 8 + "px",
        position: "absolute",
    }

    const updatePieces = (initialSquare, newSquare) => {
        const initialSquareNum = Number(initialSquare);
        const newSquareNum = Number(newSquare);

        const updatedPieces = pieces["pieces"].map(piece => {
            if (piece.offset === initialSquareNum) {
                return { ...piece, offset: newSquareNum };
            }
            return piece;
        });

        setPieces(prevPieces => ({
            ...prevPieces,
            pieces: updatedPieces
        }));
    };

    const onDrop = e => {
        const [p, initialSquare] = e.dataTransfer.getData('text').split(',');

        const {width, left, top} = ref.current.getBoundingClientRect();
        const size = width / 8;

        const rank = Math.floor((e.clientY - top) / size);
        const file = Math.floor((e.clientX - left) / size);
        const newSquare = file + rank * 8;

        updatePieces(initialSquare, newSquare);
    }

    const onDragOver = e => {
        e.preventDefault();
    }

    return (
        <div 
            style={piecesWrapperStyles}
            ref={ref}
            onDrop={onDrop}
            onDragOver={onDragOver}
        >
            {pieces.pieces.map(({ id, piece, offset }, index) => (
                <Piece piece={piece} square={offset}/>
            ))}
        </div>
    );
}



export default Pieces;