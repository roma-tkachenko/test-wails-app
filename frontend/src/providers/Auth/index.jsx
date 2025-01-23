import React, { useContext, createContext } from 'react';
import useAuth from './use';

const AuthContext = createContext({});
const useAuthContext = () => useContext(AuthContext);

const AuthContextProvider = ({ children }) => {
    const Auth = useAuth();

    return (
        <AuthContext.Provider value={Auth}>
            {children}
        </AuthContext.Provider>
    );
};

export default AuthContextProvider;
export  { useAuthContext };