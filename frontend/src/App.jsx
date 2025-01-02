import {useState, useEffect} from 'react';
import './App.css';
import {Authenticate, Greet} from "../wailsjs/go/main/App";
import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime';

function App() {
    const [authStatus, setAuthStatus] = useState(false);
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

    const eventCallback = params => {setAuthStatus(params)}

    useEffect(()=>{
        EventsOn("authStatus", eventCallback);
        return () => {EventsOff("authStatus")}
    },[])

    return (
        <>
            { !authStatus &&
                <div className="overlay">
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
            }
            { authStatus &&
                <>
                    <div id="left-sidebar">
                        <div className="user-login"></div>
                        <div className="home"></div>
                        <div className="club"></div>
                        <div className="cards"></div>
                        <div className="suppot"></div>
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
            }
        </>
    )
}

export default App
