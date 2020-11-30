import React from "react";
import ReactDOM from "react-dom";
import ReactWordcloud from "react-wordcloud";

import "tippy.js/dist/tippy.css";
import "tippy.js/animations/scale.css";

const options = {
  colors: ["#1f77b4", "#ff7f0e", "#2ca02c", "#d62728", "#9467bd", "#8c564b"],
  enableTooltip: true,
  deterministic: false,
  fontFamily: "impact",
  fontSizes: [5, 60],
  fontStyle: "normal",
  fontWeight: "normal",
  padding: 1,
  rotations: 0,
  rotationAngles: [0, 90],
  scale: "sqrt",
  spiral: "archimedean",
  transitionDuration: 1000
};

class WordCloud extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      isLoaded: false,
      wordValues: []
    };
  }

  componentDidMount() {
    fetch("https://kyle.evans.dev/words")
      .then(res => res.json())
      .then(
        (result) => {
          this.setState({
            isLoaded: true,
            wordValues: result.wordValues,
          });
        },
        // Note: it's important to handle errors here
        // instead of a catch() block so that we don't swallow
        // exceptions from actual bugs in components.
        (error) => {
          this.setState({
            isLoaded: true,
            error
          });
        }
      )
  }

  render() {
    const { error, isLoaded, wordValues } = this.state;
    if (error) {
      return <div>Error: {error.message}</div>;
    } else if (!isLoaded) {
      return <div>Loading...</div>;
    } else {
      return (
        <div>
          <div style={{ height: 400, width: 600 }}>
            <ReactWordcloud options={options} words={wordValues.wordValues} />
          </div>
        </div>
      );
    }
  }
}


const rootElement = document.getElementById("root");
ReactDOM.render(<WordCloud />, rootElement);
