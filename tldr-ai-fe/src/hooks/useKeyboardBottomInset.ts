import {useEffect, useState} from 'react';
import {Keyboard, Platform} from 'react-native';

/**
 * Keyboard overlap height for padding scroll content (iOS does not shrink the
 * root view when the keyboard opens; Android often uses adjustResize instead).
 */
export function useKeyboardBottomInset(): number {
  const [height, setHeight] = useState(0);

  useEffect(() => {
    if (Platform.OS !== 'ios') {
      return;
    }
    const show = Keyboard.addListener('keyboardWillShow', e => {
      setHeight(e.endCoordinates.height);
    });
    const hide = Keyboard.addListener('keyboardWillHide', () => {
      setHeight(0);
    });
    return () => {
      show.remove();
      hide.remove();
    };
  }, []);

  if (Platform.OS !== 'ios') {
    return 0;
  }
  return height;
}
