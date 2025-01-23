import {useEffect, useState} from "react";
import {EventsOff, EventsOn} from "../../../wailsjs/runtime/runtime.js";
import {Login} from "../../../wailsjs/go/main/App.js";


const useAuth = () => {
    const [authStatus, setAuthStatus] = useState(false);
    const authEventCallback = params => {setAuthStatus(params)}


    function login(username, password) {
        return Login(username, password);
    }

    useEffect(()=>{
        EventsOn("authStatus", authEventCallback);
        return () => {EventsOff("authStatus")}
    },[])

    return {
        authStatus,
        login
    };
};

export default useAuth;
