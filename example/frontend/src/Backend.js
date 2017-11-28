import _ from 'lodash';
import Promise from 'bluebird';
import React, {Component} from 'react';
import FlipMove from 'react-flip-move';
import TimeAgo from 'react-timeago';

import File from './File';

class Backend extends Component {

  state = {
    response: {},
    restarts: [],
    files: [],
    updates: {},
    max: 5,
    stop: false,
  }

  refreshState = () => {
    Promise.resolve(fetch(`/backend/${this.props.ip}`))
      .timeout(1000)
      .then((resp) => {
      if (resp.status !== 200) { return Promise.reject("failed") }

      return resp.json();
      })
      .then(this.updateState)
      .catch(Promise.TimeoutError, () => this.setState({stop: true}))
      .catch((err) => console.log(err))
      .finally(() => {
        if (this.state.stop) { return; }
        _.delay(this.refreshState, 1000);
      });
  }

  updateState = (body) => {
    if (_.isEqual(body, this.state.response)) {
      return;
    }

    let restarts = this.state.restarts;
    if (body.restart !== _.first(this.state.restarts)) {
      restarts = [body.restart].concat(this.state.restarts);
    }

    let files = [];
    if (body.files) {
      files = _.sortBy(body.files, (o) => o.mtime);
      _.reverse(files);
    }

    let max = this.state.max;
    const updates = _.reduce(files, (acc, {name, mtime}) => {
      const entry = _.get(this.state.updates, name, {
        count: 0,
        last: mtime
      });

      if (entry.last !== mtime) {
        entry.count += 1;
        entry.last = mtime;
        if (entry.count > max) {
          max = entry.count;
        }
      }

      acc[name] = entry;

      return acc;
    }, {});

    this.setState({
      response: body,
      restarts: restarts,
      files: files,
      updates: updates,
      max: max,
    });
  }

  componentWillMount() {
    this.refreshState();
  }

  componentWillUnmount() {
    this.setState({stop: true});
  }

  render() {
    const restarts = _.map(this.state.restarts, (restart) => (
      <li key={restart} className="list-group-item">
        <TimeAgo date={restart} />
      </li>
    ));

    const files = _.map(this.state.files, (file) => (
      <File
        {...file}
        key={file.name}
        max={this.state.max}
        updates={_.get(this.state.updates, file.name, {count: 0}).count} />
    ));

    return (
      <div className="card">
        <div className="card-header">
          <h4>Pod: {this.state.response.pod}</h4>
        </div>
        <div className="card-body">
          <div className="card-group">
            <div className="card">
              <div className="card-header bg-dark text-light">Restarts</div>
              <div className="card-body">
                <FlipMove
                  duration={750}
                  easing="ease-out"
                  typeName="ul"
                  className="list-group">
                  {restarts}
                </FlipMove>
              </div>
            </div>
            <div className="card">
              <div className="card-header bg-dark text-light">Files</div>
              <div className="card-body details-body">
                <table className="table">
                  <FlipMove duration={750} easing="ease-out" typeName="tbody">
                    {files}
                  </FlipMove>
                </table>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }
}

export default Backend;
