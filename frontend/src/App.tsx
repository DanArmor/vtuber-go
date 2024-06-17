import { useState } from 'react';
import './App.css';
import { Main } from './Pages/Main';
import { WebAppProvider } from '@vkruglikov/react-telegram-web-app';
import { Provider } from 'react-redux';
import { store } from './store/store';

function App() {
  const [smoothButtonsTransition, setSmoothButtonsTransition] = useState(false);
  return (
    <WebAppProvider options={{ smoothButtonsTransition }}>
      <Provider store={store}>
        <div className="App">
          <Main
            onChangeTransition={() => setSmoothButtonsTransition(state => !state)}
          />
        </div>
      </Provider>
    </WebAppProvider>
  );
}

export default App;
