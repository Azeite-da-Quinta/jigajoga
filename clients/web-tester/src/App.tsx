import "./App.css"
import AppContextProvider from "./components/AppContext/AppContext.tsx";
import ConnectionForm from "./components/Forms/ConnectionForm.tsx";
import InteractionInterface from "./components/InteractionInterface.tsx";

export default function App() {
    return (
        <AppContextProvider>
            <ConnectionForm />
            <InteractionInterface />
        </AppContextProvider>
    )
}