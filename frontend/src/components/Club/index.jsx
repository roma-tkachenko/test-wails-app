import React, {useEffect, useState} from 'react';
import { useAuthContext } from '../../providers/Auth'
import { BoostClub } from "../../../wailsjs/go/main/App.js";
import {EventsOff, EventsOn} from "../../../wailsjs/runtime/runtime.js";

const Club = () => {
    const { authStatus, login } = useAuthContext();

    const [boostCardUrl, setBoostCardUrl] = useState("");
    const boostCardEventCallback = params => {
        setBoostCardUrl(params)
        console.log("Event boost card")
        console.log(params)
    }

    useEffect(()=>{
        EventsOn("boostImage", boostCardEventCallback);
        return () => {EventsOff("boostImage")}
    },[])


    function boostClub() {
        BoostClub().then();
    }

    return (
        authStatus &&
        <>
            <button className="btn btn-primary btn-block btn-large" onClick={boostClub}>Boost Club</button>
            <div className="card">
                <img src={boostCardUrl} alt="Card"/>
            </div>
        </>
    );
}

export default Club;