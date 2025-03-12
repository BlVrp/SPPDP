import { StyleSheet, SafeAreaView } from 'react-native';
import Navbar from '@/components/Navbar';
import "../../global.css";
import Welcome from '@/components/ui/Welcome';

export default function HomeScreen() {
  return (
    <>
      {/* SafeAreaView ensures content appears below system display */}
      <SafeAreaView style={{ flex: 1 }}>
        <Navbar />
        <Welcome user={{
          name: 'Олеся',
          donated: 0,
          events: 0,
          raffles: 0
        }} />
      </SafeAreaView>
    </>
  );
}

const styles = StyleSheet.create({
  stepContainer: {
    gap: 8,
    marginBottom: 8,
  },
  reactLogo: {
    height: 178,
    width: 290,
    bottom: 0,
    left: 0,
    position: 'absolute',
  },
});
