export var en = {
    Phase: {
        PREGAME: '0',
        PLAY: '1'
    },
    ToServerCode: {
        DISCONNECT: '-',
        NAME: '0',
        LOBBY_CHAT_MESSAGE: '1',
        END_GAME: '2',
        DECISION: 'a',
        PROMPT_REQUEST: 'b'
    },
    ToClientCode: {
        RESTART: '0',
        LOBBY_CHAT_MESSAGE: '1',
        END_GAME: '2',
        PLAYERS: '3',
        PROMPT: 'a',
        DECISION_ACK: 'b',
        RESULT: 'c',
        WINNERS: 'd'
    }
}