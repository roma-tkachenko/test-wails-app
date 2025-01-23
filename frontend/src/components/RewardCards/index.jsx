import React, {useEffect, useState} from 'react';
import { useAuthContext } from '../../providers/Auth'
import {EventsOff, EventsOn} from "../../../wailsjs/runtime/runtime.js";
import {ClaimReward} from "../../../wailsjs/go/main/App.js";
import Card from "../Card";

const RewardCards = () => {
    const { authStatus } = useAuthContext();

    const [rewardCards, setRewardCards] = useState([]);
    const rewardCardsEventCallback = params => {
        setRewardCards(params)
        console.log("Event reward cards")
        console.log(params)
    }

    useEffect(()=>{
        EventsOn("rewardCards", rewardCardsEventCallback);
        return () => {EventsOff("rewardCards")}
    },[])

    function claimReward() {
        ClaimReward().then();
    }

    return (
        authStatus &&
        <>
            <button className="btn btn-primary btn-block btn-large" onClick={claimReward}>Get Reward Cards</button>
            <div className="cards-grid">
                {/*<div className="card">*/}
                {/*    <img src="https://animestars.org/uploads/cards_image/2437/b/formidebl-1727116358.webp"*/}
                {/*         alt="карта"/>*/}
                {/*</div>*/}
                {/*<div className="card">*/}
                {/*    <img src="https://animestars.org/uploads/cards_image/801/b/sino-asada-1731907547.webp"*/}
                {/*         alt="card"/>*/}
                {/*</div>*/}
                {/*<div className="card">*/}
                {/*    <img src="https://animestars.org/uploads/cards_image/376/b/ahjegao-dabl-ju-pis-1733409297.webp"*/}
                {/*         alt="card"/>*/}
                {/*</div>*/}
                {/*<div className="card">*/}
                {/*    <img src="https://animestars.org/uploads/cards_image/2191/b/registrator-1734467380.webp"*/}
                {/*         alt="card"/>*/}
                {/*</div>*/}
                {
                    rewardCards.map(card => (
                        <Card
                            key={card.id}
                            {...card}
                        />
                    ))
                }
            </div>
        </>
    );
}

export default RewardCards;