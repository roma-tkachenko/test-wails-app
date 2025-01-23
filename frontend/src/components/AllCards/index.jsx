import React, {useEffect, useState} from 'react';
import { useAuthContext } from '../../providers/Auth'
import {EventsOff, EventsOn} from "../../../wailsjs/runtime/runtime.js";
import {SyncCards} from "../../../wailsjs/go/main/App.js";
import Card from "../Card";

const AllCards = () => {
    const { authStatus } = useAuthContext();

    const [allCards, setSyncCards] = useState([]);
    const syncCardsEventCallback = params => {
        setSyncCards(params)
        console.log("Event reward cards")
        console.log(params)
    }

    useEffect(()=>{
        EventsOn("allCards", syncCardsEventCallback);
        return () => {EventsOff("allCards")}
    },[])

    function syncCards() {
        SyncCards().then();
    }

    return (
        authStatus &&
        <>
            <button className="btn btn-primary btn-block btn-large" onClick={syncCards}>Get All Cards</button>
            <div className="cards-grid">
                {
                    allCards.map(card => (
                        <Card
                            {...card}
                        />
                    ))
                }
            </div>
        </>
    );
}

export default AllCards;