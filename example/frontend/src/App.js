import _ from 'lodash';
import React, {Component} from 'react';

import './App.css';
import Backend from './Backend';

class App extends Component {

  state = {
    pods: [],
  }

  refreshState = () => {
    fetch('/pod-list').then((resp) => {
      if (resp.status !== 200) { return Promise.reject("failed") }

      return resp.text();
    }).then((txt) => {
      if (_.trim(txt) === "") { return; }
      this.setState({
        pods: _.split(_.trim(txt), '\n').map((v) => {
          const tmp = _.split(v, /\s+/);
          return {name: tmp[0], ip: tmp[1]};
        }),
      });
    }).catch((err) => {
      console.log(err);
    }).finally(() => {
      _.delay(this.refreshState, 5000);
    });
  }

  componentWillMount() {
    this.refreshState();
  }

  render() {
    const pods = _.map(this.state.pods, (pod) => (
      <div key={pod.name} className="row mt-5">
        <div className="col-12">
          <Backend {...pod} />
        </div>
      </div>
    ));

    return (
      <div className="container">
        {pods}
      </div>
    )
  }
}

export default App;
