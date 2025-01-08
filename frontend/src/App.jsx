import {useState, useEffect} from 'react';
import './App.css';
import { Login, ClaimReward, BoostClub } from "../wailsjs/go/main/App";
import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime';

function App() {

    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const updateName = (e) => setName(e.target.value);
    const updateResultText = (result) => setResultText(result);

    // function greet() {
    //     Greet(name).then(updateResultText);
    // }

    // Authorization logic
    const [resultLoginText, setLoginText] = useState("Please enter your credentials below 👇");
    const updateLoginText = (result) => setLoginText(result);

    const [username, setUserName] = useState('');
    const [password, setPassword] = useState('');
    const updateUserName = (e) => setUserName(e.target.value);
    const updatePassword = (e) => setPassword(e.target.value);

    const [authStatus, setAuthStatus] = useState(false);
    const authEventCallback = params => {setAuthStatus(params)}


    function login() {
        Login(username, password).then(updateLoginText);
    }

    useEffect(()=>{
        EventsOn("authStatus", authEventCallback);
        return () => {EventsOff("authStatus")}
    },[])
    // END Authorization logic

    // Reward cards logic
    const [rewardCards, setRewardCards] = useState([]);
    const rewardCardsEventCallback = params => {
        setRewardCards(params)
        console.log("Event reward cards")
        console.log(params)
    }

    function claimReward() {
        ClaimReward().then(updateResultText);
    }

    useEffect(()=>{
        EventsOn("rewardCards", rewardCardsEventCallback);
        return () => {EventsOff("rewardCards")}
    },[])
    // END Reward cards logic

    function boostClub() {
        BoostClub().then(updateResultText);
    }

    return (
        <>
            { !authStatus &&
                <>
                    <div className="overlay"></div>
                    <div className="login-form">
                        <div>
                            <input type="text" name="u" placeholder="Username" required="required"
                                   onChange={updateUserName}/>
                            <input type="password" name="p" placeholder="Password" required="required"
                                   onChange={updatePassword}/>
                            <button type="submit" className="btn btn-primary btn-block btn-large" onClick={login}>Let
                                me in.
                            </button>
                        </div>
                    </div>
                </>
            }
            { authStatus &&
                <>
                    <div id="left-sidebar">
                        {/*<div className="user-info">User Info</div>*/}
                        <ul>
                            <li className="cards"><a href="#">Cards</a></li>
                            <li className="reward-cards"><a href="#">Reward Cards</a></li>
                            <li className="boost-club"><a href="#">Boost Club</a></li>
                        </ul>
                    </div>
                    <div id="content-wrapper">
                        <div id="result" className="result">{resultText}</div>
                        <div id="input" className="input-box">
                            <input id="name" className="input" onChange={updateName} autoComplete="off" name="input"
                                   type="text"/>
                            <button className="btn" onClick={claimReward}>Get Reward Cards</button>
                            <button className="btn" onClick={boostClub}>Boost Club</button>
                        </div>
                    </div>
                </>
            }
        </>
    )
}

export default App
