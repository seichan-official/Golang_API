import React from 'react';

const HomePage = () => {
  const newUserSessionPath = '/login'; 


  return (
        <div className="container">
        <div className="row">
        <div className="mx-auto">
          <h1>
            Welcome to SpoTube!!
          </h1>
          <p>aaa</p>
          <p>aaaa</p>
          <div className="btn-wrapper col-10 mx-auto">
            <div className="row">
              <a
                href={newUserSessionPath}
                className="btn btn-info btn-sm btn-block mb-3 sign_in" >
                Log in
              </a>
            </div>
          </div>
        </div>
       </div>
      </div>
);
};


export default HomePage;
