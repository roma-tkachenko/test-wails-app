import {useState} from 'react';
import './App.css';
import {Authenticate, Greet} from "../wailsjs/go/main/App";
import { EventsOn } from '../wailsjs/runtime/runtime';

function App() {
    let countEvent = 0
    EventsOn("authStatus", (authStatus) => {
        countEvent++
        console.log("Auth status received via event:", authStatus);
        console.log("Count event:", countEvent);
    });


    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const updateName = (e) => setName(e.target.value);
    const updateResultText = (result) => setResultText(result);

    function greet() {
        Greet(name).then(updateResultText);
    }

    const [resultLoginText, setLoginText] = useState("Please enter your credentials below 👇");
    const [username, setUserName] = useState('');
    const [password, setPassword] = useState('');
    const updateLoginText = (result) => setLoginText(result);
    const updateUserName = (e) => setUserName(e.target.value);
    const updatePassword = (e) => setPassword(e.target.value);

    function authenticate() {
        Authenticate(username, password).then(updateLoginText);
    }

    return (
        <>
            <div className="overlay"></div>
            <div id="left-sidebar">
                <div className="user-login"></div>
                <div className="home"></div>
                <div className="club"></div>
                <div className="cards"></div>
                <div className="suppot"></div>
                <div className="login-form">
                    <div>
                        <div id="result" className="result">{resultLoginText}</div>
                        <input type="text" name="u" placeholder="Username" required="required"
                               onChange={updateUserName}/>
                        <input type="password" name="p" placeholder="Password" required="required"
                               onChange={updatePassword}/>
                        <button type="submit" className="btn btn-primary btn-block btn-large" onClick={authenticate}>Let
                            me in.
                        </button>
                    </div>
                </div>
            </div>
            <div id="content-wrapper">
                <div id="result" className="result">{resultText}</div>
                <div id="input" className="input-box">
                    <input id="name" className="input" onChange={updateName} autoComplete="off" name="input"
                           type="text"/>
                    <button className="btn" onClick={greet}>Greet</button>
                </div>
            </div>
        </>
    )
}

export default App
