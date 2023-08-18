import { useEffect, useState } from "react";
import logo from "./assets/images/logo-universal.png";
import "./App.css";
import { Greet, GetVersion } from "../wailsjs/go/main/App";

function App() {
    const [resultText, setResultText] = useState(
        "Please enter your name below 👇"
    );
    const [name, setName] = useState("");
    const [version, setVersion] = useState("");
    const updateName = (e: any) => setName(e.target.value);
    const updateResultText = (result: string) => setResultText(result);

    function greet() {
        Greet(name).then(updateResultText);
    }

    useEffect(() => {
        GetVersion().then((version) => setVersion(version));
    });

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo" />
            <div id="version">version: {version}</div>
            <div id="result" className="result">
                {resultText}
            </div>
            <div id="input" className="input-box">
                <input
                    id="name"
                    className="input"
                    onChange={updateName}
                    autoComplete="off"
                    name="input"
                    type="text"
                />
                <button className="btn" onClick={greet}>
                    Greet
                </button>
            </div>
        </div>
    );
}

export default App;
