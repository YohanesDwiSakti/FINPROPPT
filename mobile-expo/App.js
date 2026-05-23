import { StatusBar } from 'expo-status-bar';
import { useState } from 'react';
import {
  ActivityIndicator,
  Alert,
  SafeAreaView,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';

const API_BASE = 'http://127.0.0.1:5000';

export default function App() {
  const [activeTab, setActiveTab] = useState('track');
  const [receipt, setReceipt] = useState('');
  const [tracking, setTracking] = useState(null);
  const [origin, setOrigin] = useState('Denpasar');
  const [destination, setDestination] = useState('');
  const [weight, setWeight] = useState('');
  const [rate, setRate] = useState(null);
  const [loading, setLoading] = useState(false);

  async function trackPackage() {
    if (!receipt.trim()) {
      Alert.alert('Nomor resi kosong', 'Masukkan nomor resi dulu.');
      return;
    }

    setLoading(true);
    setTracking(null);
    try {
      const response = await fetch(`${API_BASE}/api/tracking/${encodeURIComponent(receipt.trim())}`);
      const data = await response.json();
      if (!response.ok) throw new Error(data.message || 'Gagal melacak paket');
      setTracking(data);
    } catch (error) {
      Alert.alert('Tracking gagal', error.message);
    } finally {
      setLoading(false);
    }
  }

  async function checkRate() {
    if (!origin.trim() || !destination.trim() || !weight.trim()) {
      Alert.alert('Data belum lengkap', 'Isi asal, tujuan, dan berat paket.');
      return;
    }

    setLoading(true);
    setRate(null);
    try {
      const params = new URLSearchParams({ origin, destination, weight });
      const response = await fetch(`${API_BASE}/api/rates?${params.toString()}`);
      const data = await response.json();
      if (!response.ok) throw new Error(data.message || 'Gagal menghitung ongkir');
      setRate(data);
    } catch (error) {
      Alert.alert('Cek ongkir gagal', error.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <SafeAreaView style={styles.safe}>
      <StatusBar style="light" />
      <View style={styles.header}>
        <Text style={styles.brand}>TIKI DENPASAR</Text>
        <Text style={styles.headerMeta}>Mobile Ops</Text>
      </View>

      <ScrollView contentContainerStyle={styles.content}>
        <Text style={styles.title}>Lacak paket dan cek ongkir Bali.</Text>
        <Text style={styles.subtitle}>Aplikasi Expo ini memakai backend Go yang sama dengan frontend Laravel.</Text>

        <View style={styles.tabs}>
          <TouchableOpacity
            style={[styles.tab, activeTab === 'track' && styles.tabActive]}
            onPress={() => setActiveTab('track')}>
            <Text style={[styles.tabText, activeTab === 'track' && styles.tabTextActive]}>Lacak</Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={[styles.tab, activeTab === 'rate' && styles.tabActive]}
            onPress={() => setActiveTab('rate')}>
            <Text style={[styles.tabText, activeTab === 'rate' && styles.tabTextActive]}>Ongkir</Text>
          </TouchableOpacity>
        </View>

        <View style={styles.card}>
          {activeTab === 'track' ? (
            <>
              <Text style={styles.sectionTitle}>Lacak Status Pengiriman</Text>
              <TextInput
                style={styles.input}
                value={receipt}
                onChangeText={setReceipt}
                placeholder="Masukkan nomor resi"
                autoCapitalize="characters"
              />
              <ActionButton label="Lacak Sekarang" loading={loading} onPress={trackPackage} />
              {tracking && (
                <View style={styles.result}>
                  <Text style={styles.badge}>{tracking.status}</Text>
                  <Text style={styles.resultTitle}>RESI: {tracking.receipt}</Text>
                  <Text style={styles.resultText}>Lokasi: {tracking.location}. Estimasi: {tracking.estimate}.</Text>
                  {tracking.timeline.map((step) => (
                    <View key={`${step.date}-${step.status}`} style={styles.step}>
                      <Text style={styles.stepDate}>{step.date}</Text>
                      <Text style={styles.stepStatus}>{step.status}</Text>
                    </View>
                  ))}
                </View>
              )}
            </>
          ) : (
            <>
              <Text style={styles.sectionTitle}>Hitung Estimasi Biaya</Text>
              <TextInput style={styles.input} value={origin} onChangeText={setOrigin} placeholder="Kota asal" />
              <TextInput style={styles.input} value={destination} onChangeText={setDestination} placeholder="Kota tujuan" />
              <TextInput style={styles.input} value={weight} onChangeText={setWeight} keyboardType="numeric" placeholder="Berat kg" />
              <ActionButton label="Cek Estimasi Harga" loading={loading} onPress={checkRate} />
              {rate && (
                <View style={styles.result}>
                  <Text style={styles.badge}>{rate.service}</Text>
                  <Text style={styles.resultTitle}>Rp {rate.price.toLocaleString('id-ID')}</Text>
                  <Text style={styles.resultText}>
                    {rate.origin} ke {rate.destination}, {rate.weight_kg} kg. Estimasi {rate.estimate}.
                  </Text>
                </View>
              )}
            </>
          )}
        </View>

        <View style={styles.branchCard}>
          <Text style={styles.sectionTitle}>Cabang Bali</Text>
          <Text style={styles.resultText}>Denpasar Hub, Teuku Umar, dan Sanur siap melayani pengiriman harian.</Text>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

function ActionButton({ label, loading, onPress }) {
  return (
    <TouchableOpacity style={styles.button} onPress={onPress} disabled={loading}>
      {loading ? <ActivityIndicator color="#fff" /> : <Text style={styles.buttonText}>{label}</Text>}
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1, backgroundColor: '#f7fafc' },
  header: {
    height: 76,
    backgroundColor: '#0047ff',
    paddingHorizontal: 22,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  brand: { color: '#fff', fontSize: 20, fontWeight: '800' },
  headerMeta: { color: '#dbe8ff', fontSize: 13, fontWeight: '800' },
  content: { padding: 22, paddingBottom: 40 },
  title: { color: '#0f172a', fontSize: 34, lineHeight: 38, fontWeight: '800', marginTop: 16 },
  subtitle: { color: '#64748b', fontSize: 15, lineHeight: 23, marginTop: 10 },
  tabs: { flexDirection: 'row', gap: 10, marginTop: 22, marginBottom: 14 },
  tab: { flex: 1, backgroundColor: '#e7edf5', padding: 14, borderRadius: 8, alignItems: 'center' },
  tabActive: { backgroundColor: '#0047ff' },
  tabText: { color: '#64748b', fontWeight: '800' },
  tabTextActive: { color: '#fff' },
  card: { backgroundColor: '#fff', borderRadius: 8, padding: 20, borderWidth: 1, borderColor: '#dbe4f0' },
  sectionTitle: { color: '#0f172a', fontSize: 20, fontWeight: '800', marginBottom: 14 },
  input: {
    backgroundColor: '#f7fafc',
    borderColor: '#e7edf5',
    borderWidth: 2,
    borderRadius: 8,
    padding: 15,
    marginBottom: 12,
    fontSize: 15,
  },
  button: { backgroundColor: '#0047ff', borderRadius: 8, padding: 16, alignItems: 'center', marginTop: 2 },
  buttonText: { color: '#fff', fontWeight: '800', fontSize: 15 },
  result: { marginTop: 18, backgroundColor: '#f7fafc', borderRadius: 8, borderWidth: 1, borderColor: '#dbe4f0', padding: 18 },
  badge: { color: '#059669', fontWeight: '800', fontSize: 12, marginBottom: 8 },
  resultTitle: { color: '#0f172a', fontSize: 18, fontWeight: '800', marginBottom: 6 },
  resultText: { color: '#64748b', lineHeight: 22 },
  step: { borderLeftColor: '#dbe4f0', borderLeftWidth: 3, paddingLeft: 14, marginTop: 14 },
  stepDate: { color: '#64748b', fontSize: 12, fontWeight: '800', marginBottom: 3 },
  stepStatus: { color: '#0f172a', lineHeight: 21 },
  branchCard: { marginTop: 16, backgroundColor: '#fff', borderRadius: 8, padding: 20, borderWidth: 1, borderColor: '#dbe4f0' },
});
