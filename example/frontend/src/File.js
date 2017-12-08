import moment from 'moment';
import React, {Component} from 'react';
import TimeAgo from 'react-timeago';

class File extends Component {
  render() {
    const {name, mtime, updates, max} = this.props;
    return (
      <tr className="file">
        <td>
          <div className="progress">
            <div
              className="progress-bar bg-success"
              role="progressbar"
              style={{width: (updates/max * 100) + '%'}}
              aria-valuenow="25"
              aria-valuemin="0"
              aria-valuemax="100"></div>
          </div>
        </td>
        <td>{name}</td>
        <td>
          <TimeAgo date={moment.unix(mtime)} />
        </td>
      </tr>
    )
  }
}

export default File;
