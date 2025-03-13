import { SafeAreaView } from "react-native";
import Welcome from "@/components/ui/Welcome";

export default function HomeScreen() {
  return (
    <>
      <SafeAreaView style={{ flex: 1 }}>
        <Welcome
          user={{
            name: "Олеся",
            donated: 0,
            events: 0,
            raffles: 0,
          }}
        />
      </SafeAreaView>
    </>
  );
}
