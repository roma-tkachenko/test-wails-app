import React, {useState} from 'react';
import { useAuthContext } from '../../providers/Auth'
import 'altcha';

const Authorization = () => {
    const { authStatus, login } = useAuthContext();

    const [resultLoginText, setLoginText] = useState("Please enter your credentials below 👇");
    const updateLoginText = (result) => setLoginText(result);

    const [username, setUserName] = useState('');
    const [password, setPassword] = useState('');
    const updateUserName = (e) => setUserName(e.target.value);
    const updatePassword = (e) => setPassword(e.target.value);

    function handleLogin() {
        login(username, password).then(updateLoginText);
    }

    return (
        !authStatus &&
        <>
            <div className="overlay"></div>
            <div className="login-form">
                <div>
                    <input type="text" name="u" placeholder="Username" required="required"
                           onChange={updateUserName}/>
                    <input type="password" name="p" placeholder="Password" required="required"
                           onChange={updatePassword}/>
                    <button type="submit" className="btn btn-primary btn-block btn-large" onClick={handleLogin}>Let
                        me in.
                    </button>
                    <altcha-widget
                        challengeurl="http://localhost:3000/altcha"
                        auto="onload"
                        hidefooter
                        hidelogo
                        debug
                    ></altcha-widget>
                </div>
            </div>
        </>
    );
}

export default Authorization;