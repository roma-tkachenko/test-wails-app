import {useState, useEffect} from 'react';
import './App.css';
import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime';
import { useAuthContext } from './providers/Auth'
import Authorization from './components/Authorization'
import RewardCards from './components/RewardCards'
import AllCards from "./components/AllCards";
import Club from "./components/Club";

function App() {
    const { authStatus, login } = useAuthContext();


    // const [resultText, setResultText] = useState("Please enter your name below 👇");
    // const [name, setName] = useState('');
    // const updateName = (e) => setName(e.target.value);
    // const updateResultText = (result) => setResultText(result);

    // function greet() {
    //     Greet(name).then(updateResultText);
    // }

    // Authorization logic
    // const [resultLoginText, setLoginText] = useState("Please enter your credentials below 👇");
    // const updateLoginText = (result) => setLoginText(result);

    // const [username, setUserName] = useState('');
    // const [password, setPassword] = useState('');
    // const updateUserName = (e) => setUserName(e.target.value);
    // const updatePassword = (e) => setPassword(e.target.value);
    //
    // // const [authStatus, setAuthStatus] = useState(false);
    // // const authEventCallback = params => {setAuthStatus(params)}
    //
    // function handleLogin() {
    //     login(username, password).then(updateLoginText);
    // }


    // function login() {
    //     Login(username, password).then(updateLoginText);
    // }
    //
    // useEffect(()=>{
    //     EventsOn("authStatus", authEventCallback);
    //     return () => {EventsOff("authStatus")}
    // },[])
    // END Authorization logic

    // Reward cards logic
    // const [rewardCards, setRewardCards] = useState([]);
    // const rewardCardsEventCallback = params => {
    //     setRewardCards(params)
    //     console.log("Event reward cards")
    //     console.log(params)
    // }
    //
    // useEffect(()=>{
    //     EventsOn("rewardCards", rewardCardsEventCallback);
    //     return () => {EventsOff("rewardCards")}
    // },[])
    // END Reward cards logic

    // function claimReward() {
    //     ClaimReward().then();
    // }

    const [activeBlock, setActiveBlock] = useState(1);



    return (
        <>
            <Authorization/>
            { authStatus &&
                <>
                    <div id = "left-sidebar" >
                        {/*<div className="user-info">User Info</div>*/}
                        <ul>
                            <li className="cards" onClick={() => setActiveBlock(1)}><a href="#">Cards</a></li>
                            <li className="reward-cards" onClick={() => setActiveBlock(2)}><a href="#">Reward Cards</a></li>
                            <li className="boost-club" onClick={() => setActiveBlock(3)}><a href="#">Boost Club</a></li>
                        </ul>
                    </div>
                    <div id="content-wrapper">
                        {
                            activeBlock === 1 &&
                            <AllCards/>
                        }
                        {
                            activeBlock === 2 &&
                            <RewardCards/>
                        }
                        {
                            activeBlock === 3 &&
                            <Club/>
                        }

                    </div>
                </>
            }
        </>
    )
}

export default App
