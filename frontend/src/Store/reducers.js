const initialState = {
    isUserTurn: true,
};

const appReducer = (state = initialState, action) => {
    switch (action.type) {
        case "CAN_DRAG":
            return {
                ...state,
                isUserTurn: true,
            };
        case "CANNOT_DRAG":
            return {
                ...state,
                isUserTurn: false,
            };
        default:
            return state;
    }
};

export default appReducer;
